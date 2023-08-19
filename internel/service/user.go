package service

import (
	"fmt"
	"gitlab.dian.org.cn/dianinternal/feishusign/internel/config"
	"time"

	"gitlab.dian.org.cn/dianinternal/feishusign/internel/logger"
	"gitlab.dian.org.cn/dianinternal/feishusign/internel/model"
)

const (
	MaxSignChanLen = 300
	MaxRetryTime   = 3
)

var DefaultUserService *UserService = nil
var testCount = 0

type UserService struct {
	SignMessage chan SignCode
	Exit        chan struct{}
}

type SignCode struct {
	Code      string
	Meeting   string
	RetryTime int8
}

func (u *UserService) loopDealCode() {
	for {
		select {
		case <-u.Exit:
			return
		case msg := <-u.SignMessage:
			var sign *model.SignIn

			if config.TestMode {
				testCount++
				sign = &model.SignIn{
					UserID:     fmt.Sprintf("test-%d", testCount),
					UserName:   "test",
					Status:     model.Scan,
					MeetingID:  msg.Meeting,
					CreateTime: time.Now().UnixMilli(),
				}
			} else {
				//step 0 check if out limit RetryTime
				if msg.RetryTime > MaxRetryTime {
					logger.GetLogger().Error(fmt.Sprintf("Error: out retry limit %v", msg))
					continue
				}
				//step 1 get userid and username by code
				userID, userName, err := model.GetUserMsgByCode(msg.Code)
				if err != nil {
					logger.GetLogger().Error(fmt.Sprintf("get user msg err: %v", err.Error()))
					u.SignMessage <- msg
					continue
				}
				//step 2 check if sign before
				sign, err = model.GetSignLogByIDs(userID, msg.Meeting)
				if err != nil && err != model.NotFind {
					logger.GetLogger().Error(fmt.Sprintf("get user msg err: %v", err.Error()))
					u.SignMessage <- msg
					continue
				}
				if sign.CreateTime != 0 {
					logger.GetLogger().Debug(fmt.Sprintf("DEBUG: user sign repeat %v", userName))
					continue
				}
				//step 3 insert sign log
				sign = &model.SignIn{
					UserID:     userID,
					UserName:   userName,
					Status:     model.Scan,
					MeetingID:  msg.Meeting,
					CreateTime: time.Now().UnixMilli(),
				}
			}
			if err := sign.Insert(); err != nil {
				logger.GetLogger().Error(fmt.Sprintf("insert user msg error: %v", err.Error()))
				// TODO 这里有可能由于限流导致插入不成功，需解决这种情况
				u.SignMessage <- msg
				continue
			}
			if config.TestMode {
				if err := sign.Delete(); err != nil {
					logger.GetLogger().Error(fmt.Sprintf("delete user msg error: %v", err.Error()))
				} else {
					logger.GetLogger().Debug(fmt.Sprintf("deal test %v", testCount))
				}
				continue
			}
			logger.GetLogger().Debug(fmt.Sprintf("DEBUG: user sign in success %v", sign.UserName))
		}
	}
}

func NewUserService() *UserService {
	if DefaultUserService == nil {
		DefaultUserService = &UserService{
			SignMessage: make(chan SignCode, MaxSignChanLen),
			Exit:        make(chan struct{}),
		}
		go DefaultUserService.loopDealCode()
	}
	return DefaultUserService
}
