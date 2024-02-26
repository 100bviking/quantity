package klines

import (
	"quantity/common"
	"quantity/common/db"
)

func saveKLinesPrice(kLines []*common.KLine) (err error) {
	kline := new(common.KLine)
	err = db.KLinesDB.Model(kline).Create(kLines).Error
	if err != nil {
		return
	}
	return
}
