package strategy

import (
	"quantity/common"
	"time"
)

// AvgPriceDownStrategy   代表7日均线向下穿过25日均线和99日均线
type AvgPriceDownStrategy struct {
}

func NewAvgPriceDownStrategy() Strategy {
	return &AvgPriceDownStrategy{}
}

func (down *AvgPriceDownStrategy) Analysis(symbol string, prices map[common.Interval][]*common.Price) (action *common.SubmitOrder, err error) {
	day7Prices := prices[common.Day7]
	day25Prices := prices[common.Day25]
	day99Prices := prices[common.Day99]

	// 获取当前价格
	now := time.Now()
	price, err := getCurrentPrice()
	if err != nil {
		return
	}

	currentPrice := price[symbol]
	action = &common.SubmitOrder{
		Symbol:    symbol,
		Price:     currentPrice,
		Action:    common.Hold,
		Timestamp: now,
	}

	length := len(day7Prices)
	day7LastPrice := day7Prices[0].Price
	day25LastPrice := day25Prices[0].Price
	day99LastPrice := day99Prices[0].Price

	// 当前价格下跌,7日均线价格小于25日均线价格，25日均线价格小于99日均线
	if day7LastPrice < day25LastPrice && day25LastPrice < day99LastPrice {
		// 历史一周以内发生过穿越行为,7日均线穿过了25日和99日均线
		for i := 0; i < length; i++ {
			if day7Prices[i].Price > day25Prices[i].Price && day7Prices[i].Price > day99Prices[i].Price {
				action.Action = common.Sell
				break
			}
		}
	}
	return
}
