package service

import (
	"fmt"
	"gitlab.dian.org.cn/dianinternal/feishusign/internel/config"
	"gitlab.dian.org.cn/dianinternal/feishusign/internel/tools"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"gitlab.dian.org.cn/dianinternal/feishusign/internel/logger"
	"gitlab.dian.org.cn/dianinternal/feishusign/internel/model"
)

type AdminService struct {
}

var DefaultAdminService *AdminService = nil

const (
	dataStr = "20060102"
)

func checkPrivilege(userId string) (bool, error) {
	userPart, err := model.GetUserPartByID(userId)
	if err != nil {
		return false, err
	}
	for _, r := range config.GlobalConfig.Feishu.Root {
		for _, v := range userPart {
			if v == r {
				return true, nil
			}
		}
	}
	return false, nil
}

func sentHTTPReq(method string, url string, head map[string]string, body io.Reader) ([]byte, error) {
	c := http.Client{}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	for k, v := range head {
		req.Header.Set(k, v)
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		err0 := fmt.Errorf(resp.Status)
		return nil, err0
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return respBody, nil
}

func (a *AdminService) AdminLogin(code string) (string, error) {
	//step 0 get user message
	userId, userName, err := model.GetUserMsgByCode(code)
	if err != nil {
		return "", fmt.Errorf("get user msg by code err:%v", err)
	}

	//step 1 judge user part and if root
	if ok, err := checkPrivilege(userId); !ok || err != nil {
		return "", fmt.Errorf("no privilege %v", err)
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
