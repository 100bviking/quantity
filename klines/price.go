package klines

import (
	"fmt"
	"github.com/robfig/cron"
	"quantity/common"
	"time"
)

func saveKPrice() {
	prices, err := common.FetchPrices()
	if err != nil || len(prices) == 0 {
		fmt.Println("fetch price error.")
		return
	}
	fmt.Println("save KPrice count:", len(prices))
	err = saveHistoryPrice(prices)
	if err != nil {
		fmt.Println("save KPrice error.")
		return
	}
}

func savePrice() {
	prices, err := common.FetchPrices()
	if err != nil || len(prices) == 0 {
		fmt.Println("fetch price error.")
		return
	}
	fmt.Println("save price count:", len(prices))
	err = common.SaveCurrentPrice(prices)
	if err != nil {
		fmt.Println("save current price error.")
		return
	}
}

func Run() {
	fmt.Println("start klines service.")
	c := cron.New()

	err := c.AddFunc("0 */15 * * * *", func() {
		fmt.Println("start run klines", time.Now())
		saveKPrice()
		fmt.Println("success run kines", time.Now())
	})
	if err != nil {
		panic("failed to add crontab run in klines")
	}

	err = c.AddFunc("0 * * * * *", func() {
		fmt.Println("start run price", time.Now())
		savePrice()
		fmt.Println("success run price", time.Now())
	})
	if err != nil {
		panic("failed to add crontab run in klines")
	}

	err = c.AddFunc("0 0 * * * *", func() {
		fmt.Println("start clear history price", time.Now())
		err = clearHistoryPrice()
		if err != nil {
			fmt.Println("clear history price error.")
		}
		fmt.Println("success clear history price", time.Now())
	})

	err = c.AddFunc("@weekly", func() {
		fmt.Println("start optimize price table ", time.Now())
		err = optimizePriceTable()
		if err != nil {
			fmt.Println("optimize price table  error.")
		}
		fmt.Println("success optimize price table", time.Now())
	})

	fmt.Println("successfully register klines run")
	c.Start()
	select {}
}
