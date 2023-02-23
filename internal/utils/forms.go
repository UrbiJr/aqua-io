package utils

import (
	"errors"
)

func IsStringEmpty(s string) error {
	if len(s) <= 0 {
		return errors.New("string is empty")
	}
	return nil
}
