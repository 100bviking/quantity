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

	// 当前价格上涨
	gap1 := day7LastPrice - day25LastPrice
	gap2 := day25LastPrice - day99LastPrice
	if sum == 0 &&
		gap1/day25LastPrice > 1 &&
		gap2/day99LastPrice > 1 &&
		day7LastPrice > day25LastPrice &&
		day25LastPrice > day99LastPrice &&
		gap1 > gap2 {
		// symbol当前没有买单,或者买卖单数量相同的情况下可以再次买入
		action.Action = common.Buy
	}
	return
}
