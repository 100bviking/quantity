package common

import (
	"context"
	"github.com/adshao/go-binance/v2"
	"strconv"
	"strings"
	"time"
)

const (
	CURRENT_PRICE string = "CURRENT_PRICE"
)

var (
	client *binance.Client
)

type Price struct {
	ID        int64     `gorm:"column:id"`
	Symbol    string    `gorm:"column:symbol"`
	Pair      string    `gorm:"column:pair"`
	Price     float64   `gorm:"column:price"`
	Timestamp time.Time `gorm:"column:timestamp"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (p *Price) TableName() string {
	return "price"
}

func init() {
	user, err := GetUser()
	if err != nil {
		panic("user account error")
	}
	client = binance.NewClient(user.ApiKey, user.ApiSecret)
}

func FetchPrice() (prices []*Price, err error) {
	now := time.Unix(time.Now().Unix()/60, 0)
	symbolPrices, err := client.NewListPricesService().Do(context.Background())
	if err != nil {
		return
	}

	for _, p := range symbolPrices {
		if strings.HasSuffix(p.Symbol, "USDT") {
			symbol := strings.TrimSuffix(p.Symbol, "USDT")
			price, _ := strconv.ParseFloat(p.Price, 64)
			prices = append(prices, &Price{
				Symbol:    symbol,
				Pair:      p.Symbol,
				Price:     price,
				Timestamp: now,
			})
		}
	}
	return
}