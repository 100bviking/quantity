package common

import (
	"errors"
	"gorm.io/gorm"
	"math/big"
	"strings"
)

func IngoreNotFoundError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	return err
}

func Str2Big(amount string, decimals int) (*big.Int, error) {
	amt, _, err := big.ParseFloat(amount, 10, 256, big.ToNearestEven)
	if err != nil || amt.Sign() < 0 {
		return nil, errors.New("Failed to parse amount: " + amount)
	}

	index := strings.Index(amount, ".")
	decimalsLen := 0

	if index >= 0 {
		decimalsLen = len(amount[index+1:])
		amount = amount[:index] + amount[index+1:]
	}

	if decimals > decimalsLen {
		amount = amount + strings.Repeat("0", decimals-decimalsLen)
	} else {
		amount = amount[:len(amount)-(decimalsLen-decimals)]
	}

	res, _ := big.NewInt(0).SetString(string(amount), 10)
	floatRes, _ := big.NewFloat(0).SetString(string(amount))

	base, _ := big.NewFloat(0).SetString("1" + strings.Repeat("0", decimals))
	amt = amt.Mul(amt, base)

	max := big.NewFloat(0).Mul(amt, big.NewFloat(1.1))
	min := big.NewFloat(0).Mul(amt, big.NewFloat(0.9))

	if floatRes.Cmp(max) <= 0 && floatRes.Cmp(min) >= 0 {
		return res, nil
	} else {
		bigAmt, _ := amt.Int(nil)
		return bigAmt, nil
	}
}

func Big2Str(x *big.Int, decimal int) string {
	if decimal == 0 {
		return x.String()
	}

	// Integral part
	i := x.String()
	if len(i) <= decimal {
		i = "0"
	} else {
		i = i[0 : len(i)-decimal]
	}

	// Decimal part
	d := x.String()
	if len(d) < decimal {
		d = strings.Repeat("0", decimal-len(d)) + d
	} else {
		d = d[len(d)-decimal:]
	}

	return i + "." + d
}
