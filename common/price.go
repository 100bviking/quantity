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
	ORDER        string = "ORDER"
)

const (
	Day7  Interval = 7
	Day25 Interval = 25
	Day99 Interval = 99
)

const (
	Hour = 3600
)

var (
	client          *binance.Client
	BlacklistSymbol map[string]struct{}
	WhitelistSymbol map[string]struct{}
)

type Cursor struct {
	Symbol    string    `gorm:"column:symbol"`
	Timestamp time.Time `gorm:"column:timestamp"`
	Balance   float64   `gorm:"column:balance"`
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

	BlacklistSymbol = map[string]struct{}{
		"AEUR": {},
		"ERD":  {},
		"USDC": {},
	}
	WhitelistSymbol = map[string]struct{}{
		"JUP": {},
	}
}

func FetchPrices() (prices map[string]*Price, err error) {
	now := time.Now()
	symbolPrices, err := client.NewListPricesService().Do(context.Background())
	if err != nil {
		return
	}

	prices = make(map[string]*Price)

	for _, p := range symbolPrices {
		if strings.HasSuffix(p.Symbol, "USDT") {
			symbol := strings.TrimSuffix(p.Symbol, "USDT")

			// 忽略稳定币，以及一些UP和DOWN的币
			if _, okW := WhitelistSymbol[symbol]; !okW {
				if _, okB := BlacklistSymbol[symbol]; okB {
					continue
				}
				if strings.HasSuffix(symbol, "USD") ||
					strings.HasSuffix(symbol, "UP") ||
					strings.HasSuffix(symbol, "DOWN") ||
					strings.HasSuffix(symbol, "BULL") ||
					strings.HasSuffix(symbol, "BEAR") ||
					strings.HasSuffix(symbol, "USDP") {
					continue
				}
			}

			price, _ := strconv.ParseFloat(p.Price, 64)
			prices[symbol] = &Price{
				Symbol:    symbol,
				Price:     price,
				Timestamp: now,
			}
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

func SaveCurrentPrice(prices map[string]*Price) error {
	pipeline := db.Redis.Pipeline()
	ctx := context.Background()
	for symbol, price := range prices {
		pipeline.HSet(ctx, CurrentPrice, symbol, price.Price)
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
