package strategy

import (
	"fmt"
	"quantity/common"
	"time"
)

// FifteenDownStrategy  代表15分钟线3次连续下跌提示卖出点
type FifteenDownStrategy struct {
}

func NewFifteenDownStrategy() Strategy {
	return &FifteenDownStrategy{}
}

func (f *FifteenDownStrategy) Analysis(symbol string, prices []*common.Price) (action *common.SubmitOrder, err error) {
	if len(prices) < 3 {
		err = fmt.Errorf("strategy fifteenDownStrategy not execute,because price len not enough:%d", len(prices))
		return
	}
	prices = prices[0:3]

	// 获取当前价格
	now := time.Now()
	price, err := getCurrentPrice()
	if err != nil {
		return
	}

	action = &common.SubmitOrder{
		Symbol:    symbol,
		Price:     price[symbol],
		Action:    common.Hold,
		Timestamp: now,
	}

	// 连续三连跌,卖出信号
	if prices[0].Price < prices[1].Price && prices[1].Price < prices[2].Price {
		action.Symbol = symbol
		action.Action = common.Sell
	}
	return
}
