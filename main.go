package main

import (
	"quantity/klines"
	"quantity/manage"
)

func main() {
	go klines.Run()
	go manage.Run()

	select {}
}
