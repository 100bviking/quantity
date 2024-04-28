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
	sts = append(sts, strategy.NewHammerStrategy(), strategy.NewAvgPriceDownStrategy(), strategy.NewVolumeStrategy())
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
		for _, st := range sts {
			wg.Add(1)
			channel <- 0
			go func(symbol string, st strategy.Strategy) {
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
				order, e := st.Analysis(symbol, kLines)
				if e != nil {
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
			}(symbol, st)
		}
	}
	wg.Wait()
	close(channel)
}

func Run() {
	fmt.Println("start manage service.")
	c := cron.New()

	// 1小时分析一次k线,每小时第5分钟执行分析工作
	err := c.AddFunc("0 11 * * * *", func() {
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
