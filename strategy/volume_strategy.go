package strategy

import (
	"fmt"
	"quantity/common"
	"strconv"
	"time"
)

// VolumeStrategy   交易量放量上涨
type VolumeStrategy struct {
	name string
}

func (h *VolumeStrategy) Name() string {
	return h.name
}

func NewVolumeStrategy() Strategy {
	return &VolumeStrategy{
		name: "volume",
	}
}

func (h *VolumeStrategy) Analysis(symbol string, kLines []*common.KLine) (action *common.SubmitOrder, err error) {
	fmt.Println("开始VolumeStrategy")
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

	// 判断是否异常交易量形态
	if h.isVolume(kLines) {
		action.Action = common.Buy
	}
	return
}

func (h *VolumeStrategy) isVolume(kLines []*common.KLine) (hammer bool) {
	if len(kLines) < 2 {
		return
	}

	last := kLines[0]

	// 首先必须是上涨的
	if last.EndPrice < last.StartPrice {
		return
	}

	// 交易量是上一个小时10倍
	lastVolume, _ := strconv.ParseFloat(last.VolumeTotalUsd, 64)
	secondVolume, _ := strconv.ParseFloat(kLines[1].VolumeTotalUsd, 64)
	if lastVolume/secondVolume < 10 {
		return
	}
	return true
}
