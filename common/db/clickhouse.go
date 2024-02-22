package db

import (
	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"
)

var (
	KLinesDB *gorm.DB
)

func init() {
	dsn := "clickhouse://default:123456@localhost:9000/default?dial_timeout=10s&read_timeout=20s"
	db, err := gorm.Open(clickhouse.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to init clickhouse")
	}
	KLinesDB = db
}
