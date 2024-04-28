package common

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"quantity/common/db"
	"time"
)

type OrderStatus int

const (
	Unknown OrderStatus = 0
	Success OrderStatus = 1
	Failed  OrderStatus = 2
)

type Action int

const (
	Hold Action = 0  // 不做任何操作
	Buy  Action = 1  // 买入建议
	Sell Action = -1 // 卖出建议
)

// SubmitOrder 订单
type SubmitOrder struct {
	Symbol       string    // 符号
	Price        float64   // 交易价格
	Action       Action    // 买/卖
	Timestamp    time.Time // 策略执行结束的时间点
	StrategyName string    // 作出策略的名称
}

type Order struct {
	Id           int64
	Symbol       string      `gorm:"column:symbol"`
	OrderPrice   float64     `gorm:"column:order_price"`
	SubmitPrice  float64     `gorm:"column:submit_price"`
	Amount       string      `gorm:"column:amount"`
	Money        float64     `gorm:"column:money"`
	Action       Action      `gorm:"column:action"`
	OrderTime    time.Time   `gorm:"column:order_time"`
	Status       OrderStatus `gorm:"column:status"`
	StrategyName string      `gorm:"column:strategy_name"`
	CreatedAt    time.Time   `gorm:"column:created_at"`
	UpdatedAt    time.Time   `gorm:"column:updated_at"`
}

func (o *Order) TableName() string {
	return "order"
}

func SendSubmitOrder(order *SubmitOrder) (err error) {
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

func TakeSubmitOrder(symbol string) (order *SubmitOrder, err error) {
	ctx := context.Background()

	count, err := db.Redis.LLen(ctx, symbol).Result()
	if err != nil {
		return
	}
	if count == 0 {
		return
	}

	result, err := db.Redis.BLPop(ctx, time.Minute, symbol).Result()
	if err != nil || len(result) != 2 {
		return
	}

	order = new(SubmitOrder)
	err = json.Unmarshal([]byte(result[1]), order)
	if err != nil {
		return
	}
	return
}

func SendOrder(order *Order) (err error) {
	ctx := context.Background()
	data, err := json.Marshal(order)
	if err != nil {
		return
	}

	result := db.Redis.LPush(ctx, ORDER, string(data))
	if result.Err() != nil {
		return result.Err()
	}
	return
}

func TakeAllOrder() (orders []*Order, err error) {
	ctx := context.Background()

	count, err := db.Redis.LLen(ctx, ORDER).Result()
	if err != nil {
		return
	}
	if count == 0 {
		return
	}

	pipeline := db.Redis.Pipeline()

	for i := 1; i <= int(count); i++ {
		pipeline.LPop(ctx, ORDER)
	}

	res, err := pipeline.Exec(ctx)
	if err != nil {
		return
	}

	for _, resp := range res {
		cmdData, _ := resp.(*redis.StringCmd)
		order := new(Order)
		err = json.Unmarshal([]byte(cmdData.Val()), order)
		if err != nil {
			return
		}
		orders = append(orders, order)
	}
	return
}

func SymbolOrderSumAction(symbol string) (sum int64, err error) {
	order := new(Order)
	err = db.OrderDB.Model(order).
		Select("sum(action) as sum").
		Where("symbol = ?", symbol).
		Group("symbol").
		Pluck("sum", &sum).Error
	err = IngoreNotFoundError(err)
	if err != nil {
		return
	}
	return
}

func FetchAllOrders() (ordersMap map[string][]*Order, err error) {
	order := new(Order)
	ordersMap = make(map[string][]*Order)

	var orders []*Order
	err = db.OrderDB.Model(order).Order("created_at").Find(&orders).Error
	err = IngoreNotFoundError(err)
	if err != nil {
		return
	}

	for _, o := range orders {
		ordersMap[o.Symbol] = append(ordersMap[o.Symbol], o)
	}
	return
}

func FetchSymbolBuyLastOrder(symbol string) (order *Order, err error) {
	order = new(Order)

	err = db.OrderDB.Model(order).
		Where("symbol = ? and action = ?", symbol, Buy).
		Order("created_at desc").First(order).Error
	err = IngoreNotFoundError(err)
	if err != nil {
		return
	}
	return
}
