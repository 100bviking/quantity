package main

import (
	"quantity/manage"
	"quantity/order"
)

func main() {
	go manage.Run()
	go order.Run()
	select {}
}
