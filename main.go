package main

import (
	"quantity/klines"
)

func main() {
	go klines.Run()

	select {}
}
