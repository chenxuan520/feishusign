package model

import (
	"database/sql"
	"fmt"

	"gitlab.dian.org.cn/dianinternal/feishusign/internel/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var defaultDB *gorm.DB

var NoFind error = gorm.ErrRecordNotFound

func InitMysql(dsn string) error {
	mysqlConfig := mysql.Config{
		DSN:                       dsn,
		DefaultStringSize:         256,
		SkipInitializeWithVersion: false,
	}
	if db, err := gorm.Open(mysql.New(mysqlConfig)); err != nil {
		return err
	} else {
		defaultDB = db
		return nil
	}
}

func GetMysqlDB() *gorm.DB {
	if defaultDB == nil {
		logger.GetLogger().Error("mysql database is not initialized")
		return nil
	}
	return defaultDB
}

//to create database
func CreateDatabase(dsn string, driver string, createSql string) error {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return err
	}
	defer func(db *sql.DB) {
		err = db.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(db)
	if err = db.Ping(); err != nil {
		return err
	}
	_, err = db.Exec(createSql)
	return err
}
