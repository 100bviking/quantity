package strategy

import (
	"fmt"
	"quantity/common"
	"time"
)

// HammerStrategy   锤子线
type HammerStrategy struct {
	name string
}

func (h *HammerStrategy) Name() string {
	return h.name
}

func NewHammerStrategy() Strategy {
	return &HammerStrategy{
		name: "hammer",
	}
}

func (h *HammerStrategy) Analysis(symbol string, kLines []*common.KLine) (action *common.SubmitOrder, err error) {
	fmt.Println("开始HammerStrategy")
	// 获取当前价格
	now := time.Now()
	price, err := getCurrentPrice()
	if err != nil {
		return
	}

	currentPrice := price[symbol]
	action = &common.SubmitOrder{
		Symbol:       symbol,
		Price:        currentPrice,
		Action:       common.Hold,
		Timestamp:    now,
		StrategyName: h.name,
	}

	// 判断第1根线是否满足上涨
	if !kLines[0].IsUp() {
		return
	}

	// 判断倒数第二根线是否满足锤子线形态
	if !kLines[1].IsHammer() {
		return
	}

	// 除了最后一根线最近是否满足下跌趋势
	var ks common.KLines = kLines[1:]
	if !ks.ContinueDown() {
		return
	}

	action.Action = common.Buy
	return
}
