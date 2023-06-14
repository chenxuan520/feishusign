package model

import (
	"time"
)

type Status int8

const (
	Leave Status = 1
	Scan  Status = 2
)

type SignIn struct {
	UserID     string `gorm:"column:user_id"`
	UserName   string `gorm:"column:user_name"`
	Status     Status `gorm:"column:status"`
	MeetingID  int64  `gorm:"column:meeting_id"`
	CreateTime int64  `gorm:"column:create_time"`
}

func SignInTableName() string {
	return "sign"
}

func (this *SignIn) Insert() error {
	this.CreateTime = time.Now().UnixMilli()
	err := defaultDB.Table(SignInTableName()).Create(this).Error
	return err
}

func GetSignLogByIDs(userID string, meeting int64) (*SignIn, error) {
	sign := SignIn{}
	err := defaultDB.Table(SignInTableName()).Where("user_id = ? AND meeting_id = ?", userID, meeting).Find(&sign).Error
	return &sign, err
}

func BatchSignLogByMeeting(meetindID int64) ([]*SignIn, error) {
	signs := []*SignIn{}
	err := defaultDB.Table(SignInTableName()).Where("meeting_id = ", meetindID).Find(&signs).Error
	return signs, err
}
