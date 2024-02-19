package common

import (
	"context"
	"encoding/json"
	"quantity/common/db"
	"time"
)

type Action int

const (
	Hold Action = 0 // 不做任何操作
	Buy  Action = 1 // 买入建议
	Sell Action = 2 // 卖出建议
)

// Order 订单
type Order struct {
	Symbol    string    // 符号
	Price     float64   // 交易价格
	Action    Action    // 买/卖
	Timestamp time.Time // 策略执行结束的时间点
}

func SendOrder(order *Order) (err error) {
	ctx := context.Background()
	data, err := json.Marshal(order)
	if err != nil {
		return
	}

	result := db.Redis.LPush(ctx, order.Symbol, string(data))
	if result.Err() != nil {
		return result.Err()
	}
	return
}
