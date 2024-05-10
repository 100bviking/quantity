package klines

import (
	"fmt"
	"github.com/robfig/cron"
	"quantity/common"
	"sync"
	"time"
)

var (
	wg sync.WaitGroup
)

func saveKPrice() (err error) {
	symbols, err := common.GetCurrentSymbol()
	if err != nil {
		fmt.Println("failed to get current symbol in saveKPrice.")
		return
	}

	cursorMap, err := common.GetSymbolCursor()
	if err != nil {
		fmt.Println("failed to get symbol cursor in saveKPrice.")
		return
	}

	now := time.Now()
	channel := make(chan int, 10)
	for _, symbol := range symbols {
		channel <- 0
		wg.Add(1)
		go func(symbol string) {
			defer func() {
				<-channel
				wg.Done()
			}()

			// 获取symbol对应最大的时间戳
			var (
				startTime = now.AddDate(0, 0, -7).Unix()
				endTime   = now.Unix() - 4*3600
			)

			currentTime, ok := cursorMap[symbol]
			// 如果当前cursor小于一周内则使用当前cursor,否则取一周前数据
			if ok && startTime < currentTime.Timestamp.Unix() {
				startTime = currentTime.Timestamp.Unix()
			}
			kLinePrices, e := common.QueryHistoryKLines(symbol, startTime, endTime)
			if e != nil {
				fmt.Println("failed to query history klines", symbol, startTime, endTime, e)
				return
			}
			if len(kLinePrices) > 0 {
				err = saveKLinesPrice(kLinePrices)
				if err != nil {
					fmt.Println("saveKLinesPrice failed:", err)
					return
				}
			}

			err = common.UpdateSymbolCursor(symbol, endTime)
			if err != nil {
				fmt.Println("failed to update symbol cursor", symbol, endTime)
				return
			}
		}(symbol)
	}
	wg.Wait()
	close(channel)
	return
}

func savePrice() {
	prices, err := common.FetchPrices()
	if err != nil || len(prices) == 0 {
		fmt.Println("fetch price error.", err, len(prices))
		return
	}

	fmt.Println("save price count:", len(prices))
	err = common.SaveCurrentPrice(prices)
	if err != nil {
		fmt.Println("save current price error.")
		return
	}

	ordersMap, err := common.FetchAllOrders()
	if err != nil {
		fmt.Println("failed to fetch all orders")
		return
	}

	common.CountMoney(prices, ordersMap)
}

func Run() {
	fmt.Println("start klines service.")
	c := cron.New()

	// 4 小时运行一次,保存k线
	err := c.AddFunc("0 0 */4 * * *", func() {
		fmt.Println("start run klines", time.Now())
		err := saveKPrice()
		fmt.Println("success run kines", time.Now(), err)
	})
	if err != nil {
		panic("failed to add crontab run in klines")
	}

	// 1分钟运行一次保存当前价格到redis
	err = c.AddFunc("0 * * * * *", func() {
		fmt.Println("start run price", time.Now())
		savePrice()
		fmt.Println("success run price", time.Now())
	})
	if err != nil {
		panic("failed to add crontab run in klines")
	}

	fmt.Println("successfully register klines run")
	c.Start()
	select {}
}
