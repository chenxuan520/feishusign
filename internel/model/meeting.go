package model

import "time"

type Meeting struct {
	MeetingID    string `gorm:"column:meeting_id"`
	OriginatorID string `gorm:"column:originator_id"`
	Year         int32  `gorm:"column:year"`
	Month        int32  `gorm:"column:month"`
	Day          int32  `gorm:"column:day"`
	CreateTime   int64  `gorm:"column:create_time"`
}

func MeetingTableName() string {
	return "meeting"
}

func (m *Meeting) Insert() error {
	m.CreateTime = time.Now().UnixMilli()
	err := defaultDB.Table(MeetingTableName()).Create(m).Error
	return err
}

func GetLatestMeeting() (*Meeting, error) {
	meeting := Meeting{}
	err := defaultDB.Table(MeetingTableName()).Order("create_time DESC").First(&meeting).Error
	return &meeting, err
}
