package common

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"quantity/common/db"
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

func (k *KLine) TableName() string {
	return "kline"
}

func QueryHistoryKLines(symbol string, startTime int64, endTime int64) (kLinePrices []*KLine, err error) {
	limit := int((endTime - startTime) / 3600)
	pair := symbol + "USDT"
	start := startTime * 1000
	end := endTime * 1000
	symbolPrices, err := client.NewKlinesService().
		Symbol(pair).
		StartTime(start).
		EndTime(end).
		Interval("1h").Limit(limit).Do(context.Background())

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
		cursorMap[c.Symbol] = cursor
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
