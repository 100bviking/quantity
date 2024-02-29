package manage

import (
	"quantity/common"
	"quantity/common/db"
)

func getHistoryPrice(symbol string) (priceMap map[common.Interval][]*common.Price, err error) {
	priceMap = make(map[common.Interval][]*common.Price)

	// 获取7日均线
	var (
		day7Prices  []*common.Price
		day25Prices []*common.Price
		day99Prices []*common.Price
	)
	err = db.KLinesDB.Table("kline").
		Select("symbol,k_start_time as timestamp,avg(toFloat64OrZero(volume_total_usd)) over (PARTITION BY symbol ORDER BY timestamp rows between 6 preceding and current row ) as price").
		Where("symbol = ?", symbol).
		Order("timestamp desc").
		Limit(24).
		Find(&day7Prices).Error
	err = common.IngoreNotFoundError(err)
	if err != nil || len(day7Prices) == 0 {
		return
	}
	priceMap[common.Day7] = day7Prices

	// 获取25日均线
	err = db.KLinesDB.Table("kline").
		Select("symbol,k_start_time as timestamp,avg(toFloat64OrZero(volume_total_usd)) over (PARTITION BY symbol ORDER BY timestamp rows between 24 preceding and current row ) as price").
		Where("symbol = ?", symbol).
		Order("timestamp desc").
		Limit(24).
		Find(&day25Prices).Error
	err = common.IngoreNotFoundError(err)
	if err != nil || len(day25Prices) == 0 {
		return
	}
	priceMap[common.Day25] = day25Prices

	// 获取99日均线
	err = db.KLinesDB.Table("kline").
		Select("symbol,k_start_time as timestamp,avg(toFloat64OrZero(volume_total_usd)) over (PARTITION BY symbol ORDER BY timestamp rows between 98 preceding and current row ) as price").
		Where("symbol = ?", symbol).
		Order("timestamp desc").
		Limit(24).
		Find(&day99Prices).Error
	err = common.IngoreNotFoundError(err)
	if err != nil || len(day99Prices) == 0 {
		return
	}
	priceMap[common.Day99] = day99Prices
	return
}
