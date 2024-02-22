package main

import (
	"quantity/klines"
	"quantity/order"
)

func main() {
	go klines.Run()
	//go manage.Run()
	go order.Run()
	select {}
}
