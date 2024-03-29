package db

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

var Redis *redis.Client

func init() {
	Redis = redis.NewClient(&redis.Options{
		Addr:         "127.0.0.1:6379",
		DB:           0, // use default DB
		ReadTimeout:  time.Minute,
		WriteTimeout: time.Minute,
		PoolTimeout:  time.Minute,
	})
	err := Redis.Ping(context.Background())
	if err.Err() != nil {
		panic("init redis client failed.")
	}
}
