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

func (down *AvgPriceDownStrategy) Analysis(symbol string, kLines []*common.KLine) (action *common.SubmitOrder, err error) {
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

	sum, e := common.SymbolOrderSumAction(symbol)
	if e != nil {
		return nil, e
	}

	// 当前价格下跌,7日均线价格小于25日均线价格,即刻止损卖出
	if sum >= 1 {
		//action.Action = common.Sell
	}
	return
}
