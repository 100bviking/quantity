package klines

import (
	"context"
	"quantity/common"
	"quantity/common/db"
	"time"
)

func saveHistoryPrice(prices []*common.Price) error {
	return db.KDB.Model(&common.Price{}).Save(prices).Error
}

func clearHistoryPrice() error {
	price := new(common.Price)
	yesterday := time.Now().AddDate(0, 0, -1)
	return db.KDB.Model(price).Where("created_at <= ?", yesterday).Delete(price).Error
}

func optimizePriceTable() error {
	err := db.KDB.Raw("optimize table price").Error
	return err
}

func saveCurrentPrice(prices []*common.Price) error {
	pipeline := db.Redis.Pipeline()
	ctx := context.Background()
	for _, price := range prices {
		pipeline.HSet(ctx, common.CURRENT_PRICE, price.Symbol, price.Price)
	}
	_, err := pipeline.Exec(ctx)
	return err
}