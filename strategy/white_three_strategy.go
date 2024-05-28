package strategy

import (
	"fmt"
	"quantity/common"
	"strconv"
	"time"
)

// WhiteThreeStrategy 白色三兵上攻
type WhiteThreeStrategy struct {
	name string
}

func (h *WhiteThreeStrategy) Name() string {
	return h.name
}

func NewWhiteThreeStrategy() Strategy {
	return &WhiteThreeStrategy{
		name: "white_three",
	}
}

func (h *WhiteThreeStrategy) Analysis(symbol string, kLines []*common.KLine) (action *common.SubmitOrder, err error) {
	fmt.Println("开始WhiteThreeStrategy")
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

	// 判断是否是白色三兵形态
	if !h.isWhiteThree(kLines) {
		return
	}

	action.Action = common.Buy

	return
}

func (h *WhiteThreeStrategy) isUp(kLine *common.KLine) (is bool) {
	endPrice, _ := strconv.ParseFloat(kLine.EndPrice, 64)
	startPrice, _ := strconv.ParseFloat(kLine.StartPrice, 64)

	if endPrice < startPrice {
		return
	}
	return true
}

func (h *WhiteThreeStrategy) isWhiteThree(kLines []*common.KLine) (hammer bool) {
	if len(kLines) < 3 {
		return
	}

	// 最近3根必须连续上涨,且是光头光脚
	for i := 0; i < 3; i++ {
		if !kLines[i].IsUp() {
			return
		}
		if !kLines[i].IsNoHeadOrFoot() {
			return
		}
	}

	// 最近3根必须是上涨趋势
	var kl common.KLines = kLines[0:3]
	if !kl.ContinueUp() {
		return
	}

	return true
}
