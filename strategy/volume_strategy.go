package strategy

import (
	"fmt"
	"quantity/common"
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

	// 判断是否最近整体向下
	var ks common.KLines = kLines[1:]
	if !ks.ContinueDown() {
		return
	}

	// 判断是否最近一根k线是向上
	if !kLines[0].IsUp() {
		return
	}

	// 判断是否异常交易量形态

	if !h.isVolume(kLines) {
		return
	}
	action.Action = common.Buy
	return
}

// 美元计价
func (h *VolumeStrategy) isVolume(kLines []*common.KLine) (hammer bool) {
	var ks common.KLines = kLines[1:]

	lastVolume := kLines[0].Volume()
	// 交易量是上一个小时5倍
	if lastVolume/ks.AvgVolume() < 5.0 {
		return
	}
	return true
}
