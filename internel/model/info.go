package model

import (
	"gitlab.dian.org.cn/dianinternal/feishusign/internel/logger"
	"strconv"
)

type Info struct {
	DianID    int    `gorm:"column:dian_id"`
	StudentID string `gorm:"column:student_id"`
}

func infoTableName() string {
	return "member"
}

func GetStudentIDAndDianID(name string) (string, string) {
	i := new(Info)
	err := defaultDB.Table(infoTableName()).Select([]string{"dian_id", "student_id"}).Where("name = ?", name).First(i).Error
	if err != nil {
		logger.GetLogger().Error("get studentID err: " + err.Error())
		return "", ""
	}
	if i.DianID == 0 {
		return i.StudentID, ""
	}
	return i.StudentID, strconv.Itoa(i.DianID)
}
