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

	symbol := "ALICE"
	err = saveSymbolPrice(symbol, cursorMap)
	fmt.Println(err)
}
