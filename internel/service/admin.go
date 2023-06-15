package service

import (
	"fmt"
	"time"

	"gitlab.dian.org.cn/dianinternal/feishusign/internel/logger"
	"gitlab.dian.org.cn/dianinternal/feishusign/internel/model"
)

type AdminService struct {
}

var DefaultAdminService *AdminService = nil

const dataStr = "20011203"

func (a *AdminService) AdminLogin(code string) (string, error) {
	//TODO finish it
	//step 0 get user message

	//step 1 judge user part and if root

	//step 2 create jwt

	return "", nil
}

func (a *AdminService) AdminSend(userID, text string) {
	model.RobotSendTextMsg(userID, text)
}

func (a *AdminService) AdminGetMeeting() (string, error) {
	meeting, err := model.GetLatestMeeting()
	if err != nil {
		logger.GetLogger().Error(fmt.Sprintf("Error:%s", err.Error()))
		return "", err
	}
	return meeting.MeetingID, nil
}

func (a *AdminService) AdminDealMsg(userID, text string) error {
	t, err := time.Parse(dataStr, text)
	if err != nil {
		logger.GetLogger().Error(fmt.Sprintf("Error:%s", err))
		return err
	}
	meeting := model.Meeting{
		MeetingID:    text,
		OriginatorID: userID,
		Year:         int32(t.Year()),
		Month:        int32(t.Month()),
		Day:          int32(t.Day()),
		CreateTime:   time.Now().UnixMilli(),
	}
	err = meeting.Insert()
	if err != nil {
		return err
	}
	return nil
}

func NewAdminService() *AdminService {
	if DefaultAdminService == nil {
		DefaultAdminService = &AdminService{}
	}
	return DefaultAdminService
}
