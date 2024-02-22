package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	AccountDB *gorm.DB
	KDB       *gorm.DB
	OrderDB   *gorm.DB
)

func init() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/account?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("init account db failed")
	}
	AccountDB = db

	dsnK := "root:123456@tcp(127.0.0.1:3306)/klines?charset=utf8mb4&parseTime=True&loc=Local"
	kdb, err := gorm.Open(mysql.Open(dsnK), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("init klines db failed")
	}
	KDB = kdb

	dsnOrder := "root:123456@tcp(127.0.0.1:3306)/order?charset=utf8mb4&parseTime=True&loc=Local"
	odb, err := gorm.Open(mysql.Open(dsnOrder), &gorm.Config{})
	if err != nil {
		panic("init order db failed")
	}
	OrderDB = odb
}
