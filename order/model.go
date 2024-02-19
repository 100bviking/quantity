package order

import (
	"quantity/common"
	"quantity/common/db"
)

func saveOrders(orders []*common.Order) (err error) {
	order := new(common.Order)
	err = db.OrderDB.Model(order).Save(orders).Error
	return
}
