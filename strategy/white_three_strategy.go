package strategy

import (
	"fmt"
	"quantity/common"
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

	// 判断是否是白色三兵形态
	if !h.isWhiteThree(kLines) {
		return
	}

	action.Action = common.Buy

	return
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
