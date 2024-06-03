package klines

import (
	"fmt"
	"quantity/common"
	"testing"
)

func Test_saveKPrice(t *testing.T) {
	cursorMap, err := common.GetSymbolCursor()
	if err != nil {
		fmt.Println("failed to get symbol cursor in saveKPrice.")
		return
	}

	symbols, err := common.GetCurrentSymbol()
	if err != nil {
		fmt.Println("failed to get current symbol in saveKPrice.")
		return
	}
	for _, symbol := range symbols {
		err = saveSymbolPrice(symbol, cursorMap)
		if err != nil {
			fmt.Println(err)
		}
	}
}
