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
	err = saveOrder(submitOrder, price)
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
			order, e := common.TakeOrder(symbol)
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
}

func Run() {
	fmt.Println("start order service.")
	for {
		run()
		time.Sleep(time.Second)
	}
}
