package service

import (
	"errors"
	"fmt"
	"gitlab.dian.org.cn/dianinternal/feishusign/internel/config"
	"gitlab.dian.org.cn/dianinternal/feishusign/internel/middlerware"
	"gorm.io/gorm"
	"strconv"
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
	update bool
}

var DefaultAdminService *AdminService = nil

const (
	dataStr     = "20060102"
	maxSheetReq = 3
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

	text = strings.TrimSpace(text)

	if text == "all" {
		req := SheetReq{
			userId: userID,
			date:   text,
			update: false,
		}
		select {
		case DefaultAdminService.ReqMessage <- req:
			a.AdminSend(userID, "请求成功，请稍后")
		default:
			a.AdminSend(userID, "服务繁忙，请稍后重试")
		}
		return
	} else if len(text) > 6 && text[:6] == "change" {
		sp := strings.Split(text, " ")
		if len(sp) != 3 {
			a.AdminSend(userID, "长度错误，请检查")
			return
		}

		username, err := model.GetUsernameById(userID)
		if err != nil {
			a.AdminSend(userID, "获取姓名失败，请检查")
			return
		}

		_, err = time.Parse(dataStr, sp[1])
		if err != nil {
			a.AdminSend(userID, "时间错误，请检查")
			return
		}

		status, err := strconv.Atoi(sp[2])
		if err != nil || (status != 1 && status != 2) {
			a.AdminSend(userID, "状态错误，请检查")
			return
		}

		sign := &model.SignIn{
			UserID:     userID,
			MeetingID:  sp[1],
			UserName:   username,
			Status:     model.Status(status),
			CreateTime: time.Now().UnixMilli(),
		}

		err = sign.Insert()
		if errors.Is(err, gorm.ErrDuplicatedKey) || strings.Contains(err.Error(), "Duplicate"){
			err = sign.Update()
		}
		if err != nil {
			a.AdminSend(userID, "插入错误，请检查服务器日志")
			logger.GetLogger().Error(err.Error())
			return
		}
		a.AdminSend(userID, "修改成功")
		return
	}

	t, err := time.Parse(dataStr, text)
	if err != nil {
		a.AdminSend(userID, "请输入形如 20060102 的日期")
		return
	}

	date := t.Format(dataStr)

	// 检查是否有该meeting存在

	if _, err := model.GetMeetingByID(date); err != nil {
		if err == model.NotFind {
			a.AdminSend(userID, "该会议不存在")
		} else {
			logger.GetLogger().Error(err.Error())
			a.AdminSend(userID, "查找会议错误，请查看服务器日志排错")
		}
		return
	}

	// 为更新表格预留
	update := false
	//if len(str) > 1 && str[1] == "更新" {
	//	update = true
	//}
	req := SheetReq{
		userId: userID,
		date:   date,
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
			if date == "all" {
				url, err := a.getAllSignData()
				if err != nil {
					a.AdminSend(userId, "查找会议错误，请查看服务器日志排错")
					logger.GetLogger().Error(err.Error())
					continue
				}
				a.AdminSend(userId, url)
				continue
			}

			var url string
			if meeting, err := model.GetMeetingByID(date); err != nil {
				a.AdminSend(userId, "查找会议错误，请查看服务器日志排错")
				logger.GetLogger().Error(err.Error())
				continue
			} else {
				url = meeting.Url
			}
			if url == "" {
				// 创建表格
				var err error
				url, err = model.CreateSignDataSpreadSheet(date)
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
					url, err = model.CreateSignDataSpreadSheet(date)
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

func (a *AdminService) getAllSignData() (string, error) {
	meetings, err := model.GetAllMeeting()
	if err != nil {
		return "", err
	}
	meetingNum := len(meetings)
	chatId, err := model.GetChatID()
	if err != nil {
		return "", err
	}
	members, err := model.GetUsersByChat(chatId)
	if err != nil {
		return "", err
	}

	token, url, err := model.CreateSpreadSheet("all")
	if err != nil {
		return "", err
	}

	var values [][]string
	memberNum := len(members)

deal:
	for _, m := range members {
		userId := m[0]

		parts, err := model.GetUserPartByID(userId)
		if err != nil {
			return "", fmt.Errorf("get user parts err : %v", err)
		}
		for _, part := range parts {
			if part == "导师组" {
				memberNum--
				continue deal
			}
		}

		signs, err := model.GetSignLogById(userId)
		if err != nil {
			return "", err
		}
		index := 0
		signsNum := len(*signs)
		scanCnt := 0
		leaveCnt := 0
		for _, meetingId := range meetings {
			for index < signsNum && (*signs)[index].MeetingID < meetingId {
				index++
			}
			if index >= signsNum {
				continue
			}
			sign := (*signs)[index]
			if meetingId == sign.MeetingID {
				index++
				if sign.Status == model.Scan {
					scanCnt++
				} else {
					leaveCnt++
				}
			}

		}
		values = append(values, []string{
			m[1],
			fmt.Sprintf("签到%d次", scanCnt),
			fmt.Sprintf("请假%d次", leaveCnt),
			fmt.Sprintf("缺席%d次", meetingNum-scanCnt-leaveCnt),
		})
		values = append(values, []string{})
	}
	values = append(values, []string{fmt.Sprintf("共计%d人", memberNum)})

	sheetId, err := model.GetFirstSheetId(token)
	if err != nil {
		return "", err
	}

	if err := model.InsertItem(token, sheetId, values); err != nil {
		return "", err
	}

	return url, nil
}
