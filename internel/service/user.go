package service

import (
	"fmt"
	"time"

	"gitlab.dian.org.cn/dianinternal/feishusign/internel/logger"
	"gitlab.dian.org.cn/dianinternal/feishusign/internel/model"
)

const (
	MaxSignChanLen = 300
	MaxRetryTime   = 3
)

type UserService struct {
	SignMessage chan SignCode
	Exit        chan struct{}
}

type SignCode struct {
	Code      string
	Meeting   int64
	RetryTime int8
}

func (u *UserService) loopDealCode() {
	for {
		select {
		case <-u.Exit:
			return
		case msg := <-u.SignMessage:
			//step 0 chect if out limit RetryTime
			if msg.RetryTime > MaxRetryTime {
				logger.GetLogger().Error(fmt.Sprintln("Error: out retry limit ", msg))
				continue
			}
			//step 1 get userid and user name by code
			userID, userName, err := model.GetUserMsgByCode(msg.Code)
			if err != nil {
				logger.GetLogger().Error(fmt.Sprintln("Error:get user msg ", err.Error()))
				u.SignMessage <- msg
				continue
			}
			//step 2 check if sign before
			sign, err := model.GetSignLogByIDs(userID, msg.Meeting)
			if err != model.NoFind {
				logger.GetLogger().Debug(fmt.Sprintln("DEBUG: user sign repeat ", userName))
				continue
			}
			//step 3 insert signlog
			sign = &model.SignIn{
				UserID:     userID,
				UserName:   userName,
				Status:     model.Scan,
				MeetingID:  msg.Meeting,
				CreateTime: time.Now().UnixMilli(),
			}
			err = sign.Insert()
			if err != nil {
				logger.GetLogger().Error(fmt.Sprintln("Error:insert user msg ", err.Error()))
				u.SignMessage <- msg
				continue
			}
		}
	}
}

func NewUserService() *UserService {
	service := &UserService{
		SignMessage: make(chan SignCode, MaxSignChanLen),
		Exit:        make(chan struct{}),
	}
	go service.loopDealCode()
	return service
}
