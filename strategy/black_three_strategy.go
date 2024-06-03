package strategy

import (
	"fmt"
	"quantity/common"
	"time"
)

// BlackThreeStrategy 三只乌鸦三兵下跌
type BlackThreeStrategy struct {
	name string
}

func (h *BlackThreeStrategy) Name() string {
	return h.name
}

func NewBlackThreeStrategy() Strategy {
	return &BlackThreeStrategy{
		name: "black_three",
	}
}

func (h *BlackThreeStrategy) Analysis(symbol string, kLines []*common.KLine) (action *common.SubmitOrder, err error) {
	fmt.Println("开始BlackThreeStrategy")
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

	sum, e := common.SymbolOrderSumAction(symbol)
	if e != nil || sum != 0 {
		return nil, e
	}

	// 判断是否是黑色三乌鸦形态
	if !h.isBlackThree(kLines) {
		return
	}

	action.Action = common.Sell

	return
}

func (h *BlackThreeStrategy) isBlackThree(kLines []*common.KLine) (hammer bool) {
	if len(kLines) < 3 {
		return
	}

	// 最近3根必须连续下跌,且是光头光脚
	for i := 0; i < 3; i++ {
		if kLines[i].IsUp() {
			return
		}
		if !kLines[i].IsNoHeadOrFoot() {
			return
		}
	}

	// 最近3根必须是下跌趋势
	var kl common.KLines = kLines[0:3]
	if !kl.ContinueDown() {
		return
	}

	return true
}
