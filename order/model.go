package order

import (
	"quantity/common"
	"quantity/common/db"
	"time"
)

type OrderStatus int

const (
	Unknown OrderStatus = 0
	Success OrderStatus = 1
	Failed  OrderStatus = 2
)

type Order struct {
	Id          int64
	Symbol      string        `gorm:"column:symbol"`
	OrderPrice  float64       `gorm:"column:order_price"`
	SubmitPrice float64       `gorm:"column:submit_price"`
	Amount      string        `gorm:"column:amount"`
	Money       float64       `gorm:"column:money"`
	Action      common.Action `gorm:"column:action"`
	OrderTime   time.Time     `gorm:"column:order_time"`
	Status      OrderStatus   `gorm:"column:status"`
	CreatedAt   time.Time     `gorm:"column:created_at"`
	UpdatedAt   time.Time     `gorm:"column:updated_at"`
}

func (o *Order) TableName() string {
	return "order"
}

func saveOrder(submitOrder *common.SubmitOrder, price float64) (err error) {
	order := &Order{
		Symbol:      submitOrder.Symbol,
		SubmitPrice: submitOrder.Price,
		OrderPrice:  price,
		Amount:      "100",
		Money:       0,
		Action:      submitOrder.Action,
		OrderTime:   submitOrder.Timestamp,
		Status:      Success,
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
	}
	err = db.OrderDB.Model(order).Save(order).Error
	return
}
