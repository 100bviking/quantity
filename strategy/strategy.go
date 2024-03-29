package strategy

import (
	"quantity/common"
)

// Strategy 策略分析某个币，然后得出结论结论在当前时间生成如何交易的订单
type Strategy interface {
	Analysis(symbol string, prices map[common.Interval][]*common.Price) (*common.SubmitOrder, error)
}
