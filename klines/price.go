package klines

import (
	"context"
	"fmt"
	"github.com/adshao/go-binance/v2"
	"github.com/robfig/cron"
	"strconv"
	"strings"
	"time"
)

func run() {
	user, err := GetUser()
	if err != nil {
		panic("user account error")
	}

	now := time.Unix(time.Now().Unix()/60, 0)
	client := binance.NewClient(user.ApiKey, user.ApiSecret)
	symbolPrices, err := client.NewListPricesService().Do(context.Background())
	if err != nil {
		panic("error when get data from biance api")
	}

	var prices []*Price
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
	fmt.Println("save price count:", len(prices))
	err = savePrice(prices)
	if err != nil {
		fmt.Println("save price error.")
	}

	err = saveCurrentPrice(prices)
	if err != nil {
		fmt.Println("save current price error.")
	}
}

func Run() {
	fmt.Println("start klines service.")
	c := cron.New()
	err := c.AddFunc("0 */15 * * * *", func() {
		fmt.Println("start run klines", time.Now())
		run()
		fmt.Println("success run kines", time.Now())
	})
	if err != nil {
		panic("failed to add crontab run in klines")
	}

	err = c.AddFunc("0 0 0 * * *", func() {
		fmt.Println("start clear history price", time.Now())
		err = clearHistoryPrice()
		if err != nil {
			fmt.Println("clear history price error.")
		}
		fmt.Println("success clear history price", time.Now())
	})
	fmt.Println("successfully register klines run")
	c.Start()
	select {}
}
