package service

import (
	"fmt"
	"gitlab.dian.org.cn/dianinternal/feishusign/internel/config"
	"gitlab.dian.org.cn/dianinternal/feishusign/internel/middlerware"
	"strings"
	"time"

	"gitlab.dian.org.cn/dianinternal/feishusign/internel/logger"
	"gitlab.dian.org.cn/dianinternal/feishusign/internel/model"
)

type AdminService struct {
	ReqMessage chan SheetReq
	Exit       chan struct{}
}

type SheetReq struct {
	userId string
	date   string
	url    string
	update bool
}

var DefaultAdminService *AdminService = nil

const (
	dataStr     = "20060102"
	maxSheetReq = 5
)

func checkPrivilege(userId string) bool {
	userPart, err := model.GetUserPartByID(userId)
	if err != nil {
		logger.GetLogger().Error(err.Error())
		return false
	}
	for _, r := range config.GlobalConfig.Feishu.Root {
		for _, v := range userPart {
			if v == r {
				return true
			}
		}
	}
	return false
}

func (a *AdminService) AdminLogin(code string) (string, error) {
	//step 0 get user message
	userId, _, err := model.GetUserMsgByCode(code)
	if err != nil {
		err = fmt.Errorf("get user msg by code err: %v", err)
		logger.GetLogger().Error(err.Error())
		return "", err
	}

	//step 1 judge user part and if root
	if ok := checkPrivilege(userId); !ok {
		return "", fmt.Errorf("无权限开启会议")
	}

	//step 2 create jwt
	jwt, err := middlerware.GenerateJwt(userId)
	if err != nil {
		err = fmt.Errorf("generate JWT err: %v", err)
		logger.GetLogger().Error(err.Error())
		return "", err
	}
	return jwt, nil
}

func (a *AdminService) AdminSend(userID, text string) {
	if text == "" {
		return
	}
	// 这里需要将error中的"替换成\"，否则在发消息时会出现json反序列化错误
	text = strings.Replace(text, "\"", "\\\"", -1)
	err := model.RobotSendTextMsg(userID, text)
	if err != nil {
		logger.GetLogger().Error(err.Error())
	}
	return
}

func (a *AdminService) AdminCreateMeeting(userID string) (string, error) {
	now := time.Now()
	date := now.Format(dataStr)
	meeting, err := model.GetMeetingByID(date)
	if err != nil && err != model.NotFind {
		logger.GetLogger().Error(err.Error())
		return "", err
	}
	// has existed
	if meeting.MeetingID != "" {
		return meeting.MeetingID, nil
	}
	// not exist
	meeting = &model.Meeting{
		MeetingID:    date,
		OriginatorID: userID,
		Url:          "",
		Year:         int32(now.Year()),
		Month:        int32(now.Month()),
		Day:          int32(now.Day()),
		CreateTime:   time.Now().UnixMilli(),
	}
	err = meeting.Insert()
	if err != nil {
		logger.GetLogger().Error(fmt.Sprintf("Error:%s", err.Error()))
		return "", err
	}
	return meeting.MeetingID, nil
}

func (a *AdminService) AdminDealLeave(userId, leaveTime string) error {
	t, err := time.Parse("2006-01-02 15:04:05", leaveTime)
	if err != nil {
		return err
	}

	username, err := model.GetUsernameById(userId)
	if err != nil {
		return err
	}

	date := t.Format(dataStr)

	sign := &model.SignIn{
		UserID:     userId,
		UserName:   username,
		Status:     model.Leave,
		MeetingID:  date,
		CreateTime: time.Now().UnixMilli(),
	}
	if err := sign.Insert(); err != nil {
		return err
	}
	logger.GetLogger().Debug(fmt.Sprintf("DEBUG: user leave approval deal success %s", sign.UserName))

	return nil
}

func (a *AdminService) AdminDealMsg(userID, text string) {
	if ok := checkPrivilege(userID); !ok {
		a.AdminSend(userID, "暂无权限获取签到情况表格")
		return
	}

	str := strings.Fields(text)
	text = str[0]

	t, err := time.Parse(dataStr, text)
	if err != nil {
		a.AdminSend(userID, "请输入形如 20060102 的日期")
		return
	}

	date := t.Format(dataStr)

	// 检查是否有该meeting存在
	meeting, err := model.GetMeetingByID(date)
	if err != nil {
		if err == model.NotFind {
			a.AdminSend(userID, "没有找到该会议")
		} else {
			logger.GetLogger().Error(err.Error())
			a.AdminSend(userID, "出现查找错误，请查看服务器日志排错")
		}
		return
	}

	update := false
	if len(str) > 1 && str[1] == "更新" {
		update = true
	}
	req := SheetReq{
		userId: userID,
		date:   date,
		url:    meeting.Url,
		update: update,
	}

	select {
	case DefaultAdminService.ReqMessage <- req:
		a.AdminSend(userID, "请求成功，请稍后")
	default:
		a.AdminSend(userID, "服务繁忙，请稍后重试")
	}

	return
}

func NewAdminService() *AdminService {
	if DefaultAdminService == nil {
		DefaultAdminService = &AdminService{
			ReqMessage: make(chan SheetReq, maxSheetReq),
			Exit:       make(chan struct{}),
		}
		go DefaultAdminService.loopDealReq()
	}
	return DefaultAdminService
}

func (a *AdminService) loopDealReq() {
	for {
		select {
		case <-a.Exit:
			return
		case req := <-a.ReqMessage:
			date := req.date
			userId := req.userId
			url := req.url

			if url == "" {
				// 创建表格
				var err error
				url, err = model.CreateSpreadSheet(date)
				if err != nil {
					a.AdminSend(userId, "创建表格错误，请查看服务器日志排错")
					logger.GetLogger().Error(err.Error())
					continue
				}
			} else {
				// 检查表格是否存在
				token := model.GetSpreadsheetTokenByUrl(url)

				exist, err := model.CheckSpreadSheetIfExist(token)
				if err != nil {
					a.AdminSend(userId, "查询表格信息错误，请查看服务器日志排错")
					logger.GetLogger().Error(err.Error())
					continue
				}
				if exist {
					if req.update {
						// TODO 等到飞书支持删除数据或表格后再完成此功能
						// 目前只能实现将新数据附加在表格末尾，而无法将旧数据删除
						//if err := model.UpdateContent(date, token); err != nil {
						//	a.AdminSend(userId, "更新表格错误，请查看服务器日志排错")
						//	logger.GetLogger().Error(err.Error())
						//}
					}
				} else {
					// 链接的表格已经不存在了， 需要重新创建
					url, err = model.CreateSpreadSheet(date)
					if err != nil {
						a.AdminSend(userId, "创建表格错误，请查看服务器日志排错")
						logger.GetLogger().Error(err.Error())
						continue
					}
				}
			}
			a.AdminSend(userId, url)
		}
	}
}
