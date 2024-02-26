package order

import (
	"fmt"
	"quantity/common"
	"strconv"
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
		Action:      submitOrder.Action,
		OrderTime:   submitOrder.Timestamp,
		Status:      common.Success,
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
	}

	// 每次买入100u
	if order.Action == common.Buy {
		order.Money = 100
		order.Amount = fmt.Sprintf("%f", order.Money/price)
	}

	// 卖出,从最后一笔买入获取金额和个数
	if order.Action == common.Sell {
		lastBuyOrder, e := common.FetchSymbolBuyLastOrder(order.Symbol)
		if e != nil {
			return e
		}
		order.Amount = lastBuyOrder.Amount
		amount, _ := strconv.ParseFloat(order.Amount, 64)
		order.Money = amount * order.OrderPrice
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
