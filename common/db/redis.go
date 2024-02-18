package db

import (
	"context"
	"github.com/redis/go-redis/v9"
)

var Redis *redis.Client

func init() {
	Redis = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	err := Redis.Ping(context.Background())
	if err != nil {
		panic("init redis client failed.")
	}
}
