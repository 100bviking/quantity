package strategy

import (
	"fmt"
	"quantity/common"
	"time"
)

// FifteenUpStrategy  代表15分钟线3次连续上涨提示买入点
type FifteenUpStrategy struct {
}

func NewFifteenUpStrategy() Strategy {
	return &FifteenUpStrategy{}
}

func (f *FifteenUpStrategy) Analysis(symbol string, prices []*common.Price) (action *common.Order, err error) {
	if len(prices) < 3 {
		err = fmt.Errorf("strategy fifteenUpStrategy not execute,because price len not enough:%s", len(prices))
		return
	}
	prices = prices[0:3]

	// 获取当前价格
	now := time.Now()
	price, err := getCurrentPrice()
	if err != nil {
		return
	}

	action = &common.Order{
		Symbol:    symbol,
		Price:     price[symbol],
		Action:    common.Hold,
		Timestamp: now,
	}

	// 连续三连涨,买入信号
	if prices[0].Price > prices[1].Price && prices[1].Price > prices[2].Price {
		action.Symbol = symbol
		action.Action = common.Buy
	}
	return
}
