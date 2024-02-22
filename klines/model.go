package klines

import (
	"quantity/common"
	"quantity/common/db"
	"time"
)

func clearHistoryPrice() error {
	price := new(common.Cursor)
	yesterday := time.Now().AddDate(0, 0, -1)
	return db.KDB.Model(price).Where("created_at <= ?", yesterday).Delete(price).Error
}

func optimizePriceTable() error {
	err := db.KDB.Raw("optimize table price").Error
	return err
}

func saveKLinesPrice(kLines []*common.KLine) (err error) {
	kline := new(common.KLine)
	err = db.KLinesDB.Model(kline).Create(kLines).Error
	if err != nil {
		return
	}
	return
}
