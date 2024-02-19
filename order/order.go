package order

import (
	"fmt"
	"quantity/common"
	"sync"
	"time"
)

var (
	wg sync.WaitGroup
)

func executeOrder(submitOrder *common.SubmitOrder) (err error) {
	price, err := common.FetchSymbolPrice(submitOrder.Symbol)
	if err != nil {
		return
	}

	order := &common.Order{
		Symbol:      submitOrder.Symbol,
		SubmitPrice: submitOrder.Price,
		OrderPrice:  price,
		Amount:      "100",
		Money:       0,
		Action:      submitOrder.Action,
		OrderTime:   submitOrder.Timestamp,
		Status:      common.Success,
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
	}
	err = common.SendOrder(order)
	return
}

func run() {
	// 获取当前所有symbol
	symbols, err := common.GetCurrentSymbol()
	if err != nil || len(symbols) == 0 {
		fmt.Println("get Current Symbol failed")
		return
	}

	for _, symbol := range symbols {
		wg.Add(1)
		go func(symbol string) {
			defer wg.Done()
			// 获取 order
			order, e := common.TakeSubmitOrder(symbol)
			if e != nil {
				fmt.Printf("failed to take:%s,order err:%+v\n", symbol, e)
				return
			}
			if order == nil {
				return
			}
			e = executeOrder(order)
			if e != nil {
				fmt.Printf("failed to execute order:%v\n", order)
				return
			}

		}(symbol)
	}
	wg.Wait()

	// 执行结束统一把订单入库
	orders, err := common.TakeAllOrder()
	if err != nil {
		fmt.Println("failed to take all orders")
		return
	}

	if len(orders) == 0 {
		return
	}

	err = saveOrders(orders)
	if err != nil {
		fmt.Println("failed to save orders.")
	}
}

func Run() {
	fmt.Println("start order service.")
	for {
		run()
		time.Sleep(time.Second)
	}
}
