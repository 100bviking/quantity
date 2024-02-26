package strategy

import (
	"quantity/common"
	"time"
)

// AvgPriceUpStrategy   代表7日均线向上穿过25日均线以及99日均线
type AvgPriceUpStrategy struct {
}

func NewAvgPriceUpStrategy() Strategy {
	return &AvgPriceUpStrategy{}
}

func (up *AvgPriceUpStrategy) Analysis(symbol string, prices map[common.Interval][]*common.Price) (action *common.SubmitOrder, err error) {
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

	// 至少需要满足3小时数据
	if length < 3 {
		return
	}

	// 必须连续3次上涨，确认强势
	if !(day7Prices[0].Price > day7Prices[1].Price && day7Prices[1].Price > day7Prices[2].Price) {
		return
	}

	sum, e := common.SymbolOrderSumAction(symbol)
	if e != nil {
		return nil, e
	}

	// 当前价格上涨,当前价格大于7日均线价格，并且7日均线价格大于25日均线价格，且25日均线价格大于99日均线价格
	if sum == 0 && currentPrice > day7LastPrice && day7LastPrice > day25LastPrice && day25LastPrice > day99LastPrice {
		// 历史24小时以内发生过穿越行为,7日均线曾低于25日均线，且25日均线低于99日均线
		for i := 0; i < length; i++ {
			if day7Prices[i].Price < day25Prices[i].Price && day25Prices[i].Price < day99Prices[i].Price {
				// symbol当前没有买单,或者买卖单数量相同的情况下可以再次买入
				action.Action = common.Buy
				break
			}
		}
	}
	return
}
