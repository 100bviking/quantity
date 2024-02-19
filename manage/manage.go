package manage

import (
	"fmt"
	"github.com/robfig/cron"
	"quantity/common"
	"quantity/strategy"
	"sync"
	"time"
)

var (
	wg  sync.WaitGroup
	sts []strategy.Strategy
)

func init() {
	sts = append(sts, strategy.NewFifteenUpStrategy(), strategy.NewFifteenDownStrategy())
}

func run() {
	// 获取当前所有symbol
	symbols, err := common.GetCurrentSymbol()
	if err != nil || len(symbols) == 0 {
		fmt.Println("get Current Symbol failed")
		return
	}

	// 获取历史数据
	priceMap, err := getHistoryPrice()
	if err != nil {
		fmt.Println("get history price failed")
		return
	}

	for _, symbol := range symbols {
		wg.Add(1)
		go func(symbol string) {
			defer wg.Done()
			for _, st := range sts {
				// 按顺序执行所有策略
				order, e := st.Analysis(symbol, priceMap[symbol])
				if e != nil {
					fmt.Println("execute symbol strategy failed", symbol)
					return
				}

				if order.Action == common.Hold {
					return
				}
				e = common.SendOrder(order)
				if e != nil {
					fmt.Printf("failed to send order:%+v\n", order)
					return
				}
			}
		}(symbol)
	}
	wg.Wait()
}

func Run() {
	fmt.Println("start manage service.")
	c := cron.New()

	err := c.AddFunc("30 */15 * * * *", func() {
		fmt.Println("start run manage", time.Now())
		run()
		fmt.Println("success run manage", time.Now())
	})
	if err != nil {
		panic("failed to add cron manage")
	}
	fmt.Println("successfully register manage run")

	c.Start()
	select {}
}
