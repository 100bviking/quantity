package klines

import (
	"context"
	"quantity/common/db"
	"strconv"
	"time"
)

const (
	CURRENT_PRICE string = "CURRENT_PRICE"
)

type User struct {
	ID        int64  `gorm:"column:id"`
	UserName  string `gorm:"column:user_name"`
	ApiKey    string `gorm:"column:api_key"`
	ApiSecret string `gorm:"column:api_secret"`
}

func (u *User) TableName() string {
	return "user"
}

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

func GetUser() (user *User, err error) {
	if err != nil {
		return
	}
	user = new(User)
	err = db.AccountDB.Model(user).Where("id = 1").Take(user).Error
	return
}

func saveHistoryPrice(prices []*Price) error {
	return db.KDB.Model(&Price{}).Save(prices).Error
}

func clearHistoryPrice() error {
	price := new(Price)
	yesterday := time.Now().AddDate(0, 0, -1)
	return db.KDB.Model(price).Where("created_at <= ?", yesterday).Delete(price).Error
}

func optimizePriceTable() error {
	err := db.KDB.Raw("optimize table price").Error
	return err
}

func saveCurrentPrice(prices []*Price) error {
	pipeline := db.Redis.Pipeline()
	ctx := context.Background()
	for _, price := range prices {
		pipeline.HSet(ctx, CURRENT_PRICE, price.Symbol, price.Price)
	}
	_, err := pipeline.Exec(ctx)
	return err
}

func getCurrentPrice() (prices map[string]float64, err error) {
	ctx := context.Background()
	results, err := db.Redis.HGetAll(ctx, CURRENT_PRICE).Result()
	if err != nil {
		return
	}

	prices = make(map[string]float64)
	for symbol, value := range results {
		v, _ := strconv.ParseFloat(value, 64)
		prices[symbol] = v
	}
	return
}
