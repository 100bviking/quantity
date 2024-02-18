package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	AccountDB *gorm.DB
	KDB       *gorm.DB
)

func init() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/account?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("init account db failed")
	}
	AccountDB = db

	dsnK := "root:123456@tcp(127.0.0.1:3306)/klines?charset=utf8mb4&parseTime=True&loc=Local"
	kdb, err := gorm.Open(mysql.Open(dsnK), &gorm.Config{})
	if err != nil {
		panic("init klines db failed")
	}
	KDB = kdb
}
