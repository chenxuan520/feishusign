package model

import "time"

type Meeting struct {
	MeetingID      int64  `gorm:"column:meeting_id"`
	OriginatorID   string `gorm:"column:originator_id"`
	OriginatorName string `gorm:"column:originator_name"`
	BeginTime      int64  `gorm:"column:begin_time"`
	EndTime        int64  `gorm:"column:end_time"`
	CreateTime     int64  `gorm:"column:create_time"`
}

func MeetingTableName() string {
	return "meeting"
}

func (m *Meeting) Insert() error {
	m.CreateTime = time.Now().UnixMilli()
	err := defaultDB.Table(SignInTableName()).Create(m).Error
	return err
}
