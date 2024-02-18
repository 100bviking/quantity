package klines

import (
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	"quantity/common/db"
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

func GetUser() (user *User, err error) {
	if err != nil {
		return
	}
	user = new(User)
	err = db.AccountDB.Model(user).Where("id = 1").Take(user).Error
	return
}

func savePrice(prices []*Price) error {
	return db.KDB.Save(prices).Error
}

func clearHistoryPrice() error {
	price := new(Price)
	yesterday := time.Now().AddDate(0, 0, -1)
	return db.KDB.Model(price).Where("created_at <= ?", yesterday).Delete(price).Error
}

func saveCurrentPrice(data interface{}) error {
	str, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = db.Pool.Get().Do("Set", CURRENT_PRICE, string(str))
	return err
}

func getCurrentPrice() (data []byte, err error) {
	ret, err := redis.String(db.Pool.Get().Do("Get", CURRENT_PRICE))
	if err != nil {
		return
	}
	data = []byte(ret)
	return
}
