package common

import (
	"errors"
	"gorm.io/gorm"
)

func IngoreNotFoundError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	return err
}
