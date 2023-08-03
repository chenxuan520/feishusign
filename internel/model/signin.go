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
	UserID     string `gorm:"column:user_id;index:idx_meeting_id"`
	MeetingID  string `gorm:"column:meeting_id;index:idx_meeting_id"`
	UserName   string `gorm:"column:user_name"`
	Status     Status `gorm:"column:status"`
	CreateTime int64  `gorm:"column:create_time"`
}

func (*SignIn) TableName() string {
	return "sign"
}

func SignInTableName() string {
	return "sign"
}

func (this *SignIn) Insert() error {
	this.CreateTime = time.Now().UnixMilli()
	err := defaultDB.Table(SignInTableName()).Create(this).Error
	return err
}

func GetSignLogByIDs(userID string, meeting string) (*SignIn, error) {
	sign := SignIn{}
	err := defaultDB.Table(SignInTableName()).Where("user_id = ? AND meeting_id = ?", userID, meeting).Find(&sign).Error
	return &sign, err
}

func BatchSignLogByMeeting(meetindID int64) ([]*SignIn, error) {
	signs := []*SignIn{}
	err := defaultDB.Table(SignInTableName()).Where("meeting_id = ", meetindID).Find(&signs).Error
	return signs, err
}

func GetSignStatusById(id string, meeting string) (Status, error) {
	sign := SignIn{}
	err := defaultDB.Table(SignInTableName()).Where("user_id = ? AND meeting_id = ?", id, meeting).First(&sign).Error
	return sign.Status, err
}
