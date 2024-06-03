package common

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"math"
	"quantity/common/db"
	"strconv"
	"time"
)

type KLine struct {
	Symbol           string    `gorm:"column:symbol"`
	KStartTime       time.Time `gorm:"column:k_start_time"`
	KEndTime         time.Time `gorm:"column:k_end_time"`
	StartPrice       string    `gorm:"column:start_price"`
	EndPrice         string    `gorm:"column:end_price"`
	HighPrice        string    `gorm:"column:high_price"`
	LowPrice         string    `gorm:"column:low_price"`
	VolumeTotalUsd   string    `gorm:"column:volume_total_usd"`
	VolumeTotalCount int64     `gorm:"column:volume_total_count"`
}

// KLines 注意k线默认是逆序排列的，最新的在最前面
type KLines []*KLine

func (k *KLine) TableName() string {
	return "kline"
}

// IsNoHeadOrFoot 判断k线是否是光头光脚形态
func (k *KLine) IsNoHeadOrFoot() bool {
	endPrice, _ := strconv.ParseFloat(k.EndPrice, 64)
	startPrice, _ := strconv.ParseFloat(k.StartPrice, 64)
	lowPrice, _ := strconv.ParseFloat(k.LowPrice, 64)
	highPrice, _ := strconv.ParseFloat(k.HighPrice, 64)

	// 计算实体长度
	height := math.Abs(endPrice - startPrice)

	var (
		downHeight, upHeight float64
	)

	// 如果上涨
	if endPrice > startPrice {
		// 计算上影线
		upHeight = math.Abs(highPrice - endPrice)
		// 计算下影线
		downHeight = math.Abs(startPrice - lowPrice)
	} else {
		// 计算上影线
		upHeight = math.Abs(highPrice - startPrice)
		// 计算下影线
		downHeight = math.Abs(endPrice - lowPrice)
	}

	// 上下影线高度不能超过实体的1/5
	return (upHeight/height) < 0.2 && (downHeight/height) < 0.2
}

// IsUp 判断K线是否是上涨形态
func (k *KLine) IsUp() bool {
	endPrice, _ := strconv.ParseFloat(k.EndPrice, 64)
	startPrice, _ := strconv.ParseFloat(k.StartPrice, 64)

	return endPrice >= startPrice
}

func (k *KLine) IsHammer() bool {
	endPrice, _ := strconv.ParseFloat(k.EndPrice, 64)
	startPrice, _ := strconv.ParseFloat(k.StartPrice, 64)
	lowPrice, _ := strconv.ParseFloat(k.LowPrice, 64)

	// 计算实体长度
	height := endPrice - startPrice

	// 计算下影线
	downHeight := startPrice - lowPrice

	// 首先必须是上涨的
	// 下影线高度是实体2倍以上
	return (endPrice > startPrice) && (downHeight/height >= 2)
}

// AvgPrice 获取k线平均价
func (k *KLine) AvgPrice() float64 {
	endPrice, _ := strconv.ParseFloat(k.EndPrice, 64)
	startPrice, _ := strconv.ParseFloat(k.StartPrice, 64)
	return (endPrice + startPrice) / 2
}

// ClosePrice 获取k线收盘价
func (k *KLine) ClosePrice() float64 {
	endPrice, _ := strconv.ParseFloat(k.EndPrice, 64)
	return endPrice
}

// OpenPrice 获取k线开盘价
func (k *KLine) OpenPrice() float64 {
	startPrice, _ := strconv.ParseFloat(k.StartPrice, 64)
	return startPrice
}

func (k *KLine) Volume() float64 {
	volume, _ := strconv.ParseFloat(k.VolumeTotalUsd, 64)
	return volume
}

func (ks KLines) ContinueUp() bool {
	for i := 0; i < len(ks)-1; i++ {
		if !ks[i].IsUp() {
			return false
		}
		if ks[i].OpenPrice() < ks[i+1].ClosePrice() {
			return false
		}
	}
	return true
}

func (ks KLines) ContinueDown() bool {
	for i := 0; i < len(ks)-1; i++ {
		if ks[i].IsUp() {
			return false
		}
		if ks[i].OpenPrice() > ks[i+1].ClosePrice() {
			return false
		}
	}
	return true
}

func (ks KLines) AvgVolume() float64 {
	avgVolume := 0.0
	for i := 1; i < len(ks); i++ {
		avgVolume += ks[i].Volume()
	}
	avgVolume = avgVolume / float64(len(ks))
	return avgVolume
}

func QueryHistoryKLines(symbol string, startTime int64, endTime int64) (kLinePrices []*KLine, err error) {
	pair := fmt.Sprintf("%s%s", symbol, "USDT")
	start := startTime * 1000
	end := endTime * 1000
	symbolPrices, err := client.NewKlinesService().
		Symbol(pair).
		StartTime(start).
		EndTime(end).
		Interval("4h").Do(context.Background())
	if err != nil {
		return
	}

	for _, price := range symbolPrices {
		startAt := time.Unix(price.OpenTime/1000, 0)
		endAt := time.Unix(price.CloseTime/1000, 0)

		kLinePrices = append(kLinePrices, &KLine{
			Symbol:           symbol,
			KStartTime:       startAt,
			KEndTime:         endAt,
			StartPrice:       price.Open,
			EndPrice:         price.Close,
			HighPrice:        price.High,
			LowPrice:         price.Low,
			VolumeTotalUsd:   price.QuoteAssetVolume,
			VolumeTotalCount: price.TradeNum,
		})
	}
	return
}

func GetSymbolCursor() (cursorMap map[string]*Cursor, err error) {
	cursor := new(Cursor)
	cursorMap = make(map[string]*Cursor)
	var cursors []*Cursor
	err = db.KDB.Model(cursor).Find(&cursors).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = nil
		}
		return
	}
	for _, c := range cursors {
		cursorMap[c.Symbol] = c
	}
	return
}

func UpdateSymbolCursor(symbol string, timestamp int64) (err error) {
	cursor := &Cursor{
		Symbol:    symbol,
		Timestamp: time.Unix(timestamp, 0),
	}

	err = db.KDB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "symbol"}},
		UpdateAll: true,
	}).Create(cursor).Error
	return
}
