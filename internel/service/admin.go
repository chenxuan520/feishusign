package service

import (
	"fmt"
	"gitlab.dian.org.cn/dianinternal/feishusign/internel/config"
	"gitlab.dian.org.cn/dianinternal/feishusign/internel/tools"
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
	userId, userName, err := model.GetUserMsgByCode(code)
	if err != nil {
		return "", fmt.Errorf("get user msg by code err:%v", err)
	}

	//step 1 judge user part and if root
	if ok:= checkPrivilege(userId); !ok{
		return "", fmt.Errorf("无权限开启会议")
	}

	//step 2 create jwt
	jwt, err := tools.GenerateJwtToken(userId, userName)
	if err != nil {
		return "", err
	}
	return jwt, nil
}

func (a *AdminService) AdminSend(userID, text string) error {
	err := model.RobotSendTextMsg(userID, text)
	if err != nil {
		logger.GetLogger().Error(err.Error())
	}
	return nil
}

func (a *AdminService) AdminCreateMeeting(userID string) (string, error) {
	now := time.Now()
	date := now.Format(dataStr)
	meeting, err := model.GetMeetingByID(date)
	if err != nil && err != model.NoFind {
		logger.GetLogger().Error(err.Error())
		return "", nil
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

func (a *AdminService) AdminDealMsg(userID, text string) {
	_, err := time.Parse(dataStr, text)
	if err != nil {
		logger.GetLogger().Error(fmt.Sprintf("Error:%s", err))
		// 这里需要将error中的"进行替换，否则在发消息时会出现json反序列化错误
		// 同理，如果发送的消息中含有{}，也需要进行替换
		a.AdminSend(userID, fmt.Sprintf("error: %s", strings.Replace(err.Error(), "\"", "'", -1)))
		return
	}

	if ok:= checkPrivilege(userID); !ok {
		a.AdminSend(userID, "暂无权限获取签到情况表格")
		return
	}

	// 检查是否有该meeting存在
	meeting, err := model.GetMeetingByID(text)
	if err != nil {
		a.AdminSend(userID, err.Error())
		return
	}

	req := SheetReq{
		userId: userID,
		date:   text,
		url:    meeting.Url,
	}

	select {
	case DefaultAdminService.ReqMessage <- req:
		a.AdminSend(userID, "请求成功，请稍后")
	default:
		a.AdminSend(userID, "请求失败，触发限流")
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
					a.AdminSend(userId, fmt.Sprintln("create sheet err :", err.Error()))
					continue
				}
			} else {
				// 检查表格是否存在
				exist, err := model.CheckSpreadSheetIfExist(url)
				if err != nil {
					a.AdminSend(userId, fmt.Sprintln("create sheet err :", err.Error()))
					continue
				}
				if exist {
					// TODO 更新表格 (当然也可以不更新就是了)
				} else {
					// 链接的表格已经不存在了， 需要重新创建
					url, err = model.CreateSpreadSheet(date)
					if err != nil {
						a.AdminSend(userId, fmt.Sprintln("create sheet err :", err.Error()))
						continue
					}
				}
			}
			a.AdminSend(userId, url)
		}
	}
}
