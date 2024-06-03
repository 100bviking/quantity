package manage

import (
	"fmt"
	"github.com/robfig/cron"
	"quantity/common"
	"quantity/strategy"
	"runtime"
	"sync"
	"time"
)

var (
	wg  sync.WaitGroup
	sts []strategy.Strategy
)

func init() {
	sts = append(sts, strategy.NewHammerStrategy(),
		strategy.NewAvgPriceDownStrategy(),
		strategy.NewVolumeStrategy(),
		strategy.NewWhiteThreeStrategy(),
		strategy.NewBlackThreeStrategy(),
	)
}

func run() {
	// 获取当前所有symbol
	symbols, err := common.GetCurrentSymbol()
	if err != nil || len(symbols) == 0 {
		fmt.Println("get Current Symbol failed")
		return
	}

	channel := make(chan int, runtime.NumCPU())
	for _, symbol := range symbols {
		// 执行所有策略
		for _, singleStrategy := range sts {
			wg.Add(1)
			channel <- 0
			go func(symbol string, currentStrategy strategy.Strategy) {
				defer func() {
					wg.Done()
					<-channel
				}()

				fmt.Println("====>starting analysis symbol kline", symbol)
				// 获取历史数据
				kLines, err := getHistoryPrice(symbol)
				if err != nil || len(kLines) == 0 {
					fmt.Println("get history price failed", symbol)
					return
				}
				order, e := currentStrategy.Analysis(symbol, kLines)
				if e != nil || order == nil {
					fmt.Println("execute symbol strategy failed", symbol)
					return
				}

				if order.Action == common.Hold {
					return
				}
				e = common.SendSubmitOrder(order)
				if e != nil {
					fmt.Printf("failed to send order:%+v\n", order)
					return
				}
				fmt.Println("====>end analysis symbol kline", symbol)
			}(symbol, singleStrategy)
		}
	}
	wg.Wait()
	close(channel)
}

func Run() {
	fmt.Println("start manage service.")
	c := cron.New()

	// 4小时分析一次k线,每4小时第2分钟执行分析工作
	err := c.AddFunc("0 2 */4 * * *", func() {
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
