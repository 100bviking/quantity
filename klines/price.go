package klines

import (
	"fmt"
	"github.com/robfig/cron"
	"quantity/common"
	"sync"
	"time"
)

const (
	kLineStartTimestamp int64 = 1672502400 // 2023-01-01 00:00:00

	Hour  = 3600
	Day   = Hour * 24
	Month = Day * 30
	Week  = Day * 7
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

	now := time.Now().Unix()
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
				startTime int64
				endTime   int64
			)

			currentTime, ok := cursorMap[symbol]
			if !ok {
				startTime = kLineStartTimestamp
			} else {
				startTime = currentTime.Timestamp.Unix()
			}

			if startTime < now-Month {
				endTime = startTime + Month
			} else if startTime < now-Week {
				endTime = startTime + Week
			} else if startTime < now-Day {
				endTime = startTime + Day
			} else if startTime < now-Hour {
				endTime = startTime + Hour
			} else {
				return
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
		fmt.Println("fetch price error.")
		return
	}
	fmt.Println("save price count:", len(prices))
	err = common.SaveCurrentPrice(prices)
	if err != nil {
		fmt.Println("save current price error.")
		return
	}
}

func Run() {
	fmt.Println("start klines service.")
	c := cron.New()

	// 1 分钟运行一次
	err := c.AddFunc("0 * * * * *", func() {
		fmt.Println("start run klines", time.Now())
		err := saveKPrice()
		fmt.Println("success run kines", time.Now(), err)
	})
	if err != nil {
		panic("failed to add crontab run in klines")
	}

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
