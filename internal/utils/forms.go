package utils

import (
	"errors"
	"strconv"
)

func IsStringEmpty(s string) error {
	if len(s) <= 0 {
		return errors.New("string is empty")
	}
	return nil
}

func IsFloat(s string) error {
	if _, err := strconv.ParseFloat(s, 64); err != nil || s == "" {
		return errors.New("float conversion failed")
	}
	return nil
}

func IsInteger(s string) error {
	if _, err := strconv.ParseInt(s, 10, 64); err != nil || s == "" {
		return errors.New("integer conversion failed")
	}
	return nil
}
