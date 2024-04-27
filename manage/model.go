package manage

import (
	"quantity/common"
	"quantity/common/db"
)

func getHistoryPrice(symbol string) (kLines []*common.KLine, err error) {
	err = db.KLinesDB.Table("kline").
		Where("symbol = ?", symbol).
		Order("timestamp desc").
		Limit(7).
		Find(&kLines).Error
	err = common.IngoreNotFoundError(err)
	return
}
