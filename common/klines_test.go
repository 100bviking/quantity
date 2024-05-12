package common

import (
	"fmt"
	"testing"
	"time"
)

func TestQueryHistoryKLines(t *testing.T) {

	symbol := "BTC"
	var startTime int64 = 1714541249
	endTime := (time.Now().Unix()/(4*Hour) - 1) * (4 * Hour)
	gotKLinePrices, err := QueryHistoryKLines(symbol, startTime, endTime)
	fmt.Println("gotKLinePrices", gotKLinePrices, "err:", err)
}
