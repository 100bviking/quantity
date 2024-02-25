package common

import (
	"context"
	"github.com/adshao/go-binance/v2"
	"quantity/common/db"
	"strconv"
	"strings"
	"time"
)

type Interval int

// redis key
const (
	CurrentPrice string = "CURRENT_PRICE"
	ORDER               = "ORDER"
)

const (
	Day7  Interval = 7
	Day25 Interval = 25
	Day99 Interval = 99
)

var (
	client *binance.Client
)

type Cursor struct {
	Symbol    string    `gorm:"column:symbol"`
	Timestamp time.Time `gorm:"column:timestamp"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (p *Cursor) TableName() string {
	return "cursor"
}

type Price struct {
	Symbol    string
	Price     float64
	Timestamp time.Time
}

func init() {
	user, err := GetUser()
	if err != nil {
		panic("user account error")
	}
	client = binance.NewClient(user.ApiKey, user.ApiSecret)
}

func FetchPrices() (prices []*Price, err error) {
	now := time.Now()
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
				Price:     price,
				Timestamp: now,
			})
		}
	}
	return
}

func GetCurrentSymbol() (symbols []string, err error) {
	ctx := context.Background()
	symbols, err = db.Redis.HKeys(ctx, CurrentPrice).Result()
	if err != nil {
		return
	}
	return
}

func SaveCurrentPrice(prices []*Price) error {
	pipeline := db.Redis.Pipeline()
	ctx := context.Background()
	for _, price := range prices {
		pipeline.HSet(ctx, CurrentPrice, price.Symbol, price.Price)
	}
	_, err := pipeline.Exec(ctx)
	return err
}

func FetchSymbolPrice(symbol string) (price float64, err error) {
	ctx := context.Background()
	ret, err := db.Redis.HGet(ctx, CurrentPrice, symbol).Result()
	if err != nil {
		return
	}
	price, err = strconv.ParseFloat(ret, 64)
	return
}
