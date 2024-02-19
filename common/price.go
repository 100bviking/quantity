package common

import (
	"context"
	"github.com/adshao/go-binance/v2"
	"quantity/common/db"
	"strconv"
	"strings"
	"time"
)

// redis key
const (
	CURRENT_PRICE string = "CURRENT_PRICE"
	ORDER                = "ORDER"
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

func FetchPrices() (prices []*Price, err error) {
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

func GetCurrentSymbol() (symbols []string, err error) {
	ctx := context.Background()
	symbols, err = db.Redis.HKeys(ctx, CURRENT_PRICE).Result()
	if err != nil {
		return
	}
	return
}

func SaveCurrentPrice(prices []*Price) error {
	pipeline := db.Redis.Pipeline()
	ctx := context.Background()
	for _, price := range prices {
		pipeline.HSet(ctx, CURRENT_PRICE, price.Symbol, price.Price)
	}
	_, err := pipeline.Exec(ctx)
	return err
}

func FetchSymbolPrice(symbol string) (price float64, err error) {
	ctx := context.Background()
	ret, err := db.Redis.HGet(ctx, CURRENT_PRICE, symbol).Result()
	if err != nil {
		return
	}
	price, err = strconv.ParseFloat(ret, 64)
	return
}
