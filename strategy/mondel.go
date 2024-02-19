package strategy

import (
	"context"
	"quantity/common"
	"quantity/common/db"
	"strconv"
)

func getCurrentPrice() (prices map[string]float64, err error) {
	ctx := context.Background()
	results, err := db.Redis.HGetAll(ctx, common.CURRENT_PRICE).Result()
	if err != nil {
		return
	}

	prices = make(map[string]float64)
	for symbol, value := range results {
		v, _ := strconv.ParseFloat(value, 64)
		prices[symbol] = v
	}
	return
}
