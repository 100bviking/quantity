package strategy

import (
	"fmt"
	"quantity/common"
	"strconv"
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

	sum, e := common.SymbolOrderSumAction(symbol)
	if e != nil || sum != 0 {
		return nil, e
	}

	// 判断倒数第二根线是否满足锤子线形态
	if !h.isHammer(kLines[:len(kLines)-1]) {
		return
	}

	// 判断倒数第1根线是否满足上涨
	if !h.isUp(kLines[len(kLines)-1]) {
		return
	}
	action.Action = common.Buy

	return
}

func (h *HammerStrategy) isUp(kLine *common.KLine) (is bool) {
	endPrice, _ := strconv.ParseFloat(kLine.EndPrice, 64)
	startPrice, _ := strconv.ParseFloat(kLine.StartPrice, 64)

	if endPrice < startPrice {
		return
	}
	return true
}

func (h *HammerStrategy) isHammer(kLines []*common.KLine) (hammer bool) {
	if len(kLines) < 2 {
		return
	}

	last := kLines[0]

	endPrice, _ := strconv.ParseFloat(last.EndPrice, 64)
	startPrice, _ := strconv.ParseFloat(last.StartPrice, 64)
	lowPrice, _ := strconv.ParseFloat(last.LowPrice, 64)

	// 首先必须是上涨的
	if endPrice < startPrice {
		return
	}

	// 计算实体长度
	height := endPrice - startPrice

	// 计算下影线
	downHeight := startPrice - lowPrice

	// 下影线高度是实体2倍以上
	if downHeight/height < 2 {
		return
	}

	// 最近是下跌趋势
	firstPrice, _ := strconv.ParseFloat(kLines[len(kLines)-1].EndPrice, 64)
	if endPrice > firstPrice {
		return
	}

	return true
}
