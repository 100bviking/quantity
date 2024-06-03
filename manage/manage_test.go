package manage

import (
	"fmt"
	"quantity/strategy"
	"testing"
)

func Test_WhiteThree(t *testing.T) {
	symbol := "ALICE"
	// 获取历史数据
	kLines, err := getHistoryPrice(symbol)
	if err != nil || len(kLines) == 0 {
		fmt.Println("get history price failed", symbol)
		return
	}

	currentStrategy := strategy.NewWhiteThreeStrategy()
	order, e := currentStrategy.Analysis(symbol, kLines)
	if e != nil || order == nil {
		fmt.Println("execute symbol strategy failed", symbol)
		return
	}
}
