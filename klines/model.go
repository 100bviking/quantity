package klines

import (
	"quantity/common"
	"quantity/common/db"
)

func saveKLinesPrice(kLines []*common.KLine) (err error) {
	for _, line := range kLines {
		// check exist or not.
		var cnt int64
		err = db.KLinesDB.Model(line).Where("symbol = ? and k_start_time = ?", line.Symbol, line.KStartTime).Count(&cnt).Error
		if err != nil {
			return
		}

		if cnt > 0 {
			continue
		}
		err = db.KLinesDB.Model(line).Create(line).Error
		if err != nil {
			return
		}
	}
	return
}
