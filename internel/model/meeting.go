package model

import "time"

type Meeting struct {
	MeetingID      string `gorm:"column:meeting_id"`
	OriginatorID   string `gorm:"column:originator_id"`
	OriginatorName string `gorm:"column:originator_name"`
	Year           int32  `gorm:"column:year"`
	Month          int32  `gorm:"column:month"`
	Day            int32  `gorm:"column:day"`
	Date           string `gorm:"column:date"`
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
