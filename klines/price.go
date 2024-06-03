package klines

import (
	"fmt"
	"quantity/common"
	"sync"
	"time"
)

var (
	wg sync.WaitGroup
)

func saveSymbolPrice(symbol string, cursorMap map[string]*common.Cursor) (err error) {
	now := time.Now()
	// 获取symbol对应最大的时间戳
	var (
		startTime = now.AddDate(0, 0, -7).Unix()
		endTime   = (now.Unix()/(4*common.Hour) - 1) * (4 * common.Hour)
	)

	currentTime, ok := cursorMap[symbol]
	if ok && startTime < currentTime.Timestamp.Unix() {
		startTime = currentTime.Timestamp.Unix()
	}
	if startTime >= endTime {
		fmt.Println("开始结束时间相同，跳过")
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
			fmt.Println("saveKLinesPrice failed:", err, symbol)
			return
		}
	}

	err = common.UpdateSymbolCursor(symbol, endTime)
	if err != nil {
		fmt.Println("failed to update symbol cursor", symbol, endTime)
		return
	}
	return
}

func SaveKPrice() (err error) {
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

	channel := make(chan int, 10)
	for _, symbol := range symbols {
		channel <- 0
		wg.Add(1)
		go func(symbol string, cursorMap map[string]*common.Cursor) {
			defer wg.Done()
			saveSymbolPrice(symbol, cursorMap)
		}(symbol, cursorMap)

	}
	wg.Wait()
	close(channel)
	return
}

func SaveCurrentPrice() {
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
}
