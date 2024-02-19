package manage

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"quantity/common"
	"quantity/common/db"
)

func getCurrentSymbol() (symbols []string, err error) {
	ctx := context.Background()
	symbols, err = db.Redis.HKeys(ctx, common.CURRENT_PRICE).Result()
	if err != nil {
		return
	}
	return
}

func getHistoryPrice() (priceMap map[string][]*common.Price, err error) {
	var (
		price  = new(common.Price)
		prices = make([]*common.Price, 0)
	)
	priceMap = make(map[string][]*common.Price)
	err = db.KDB.Model(price).Order("created_at desc").
		Find(&prices).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = nil
		}
		return
	}
	for _, p := range prices {
		priceMap[p.Symbol] = append(priceMap[p.Symbol], p)
	}
	return
}