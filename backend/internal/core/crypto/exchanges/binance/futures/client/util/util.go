package util

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
)

func Sign(secretAPI, message string) string {
	mac := hmac.New(sha256.New, []byte(secretAPI))
	mac.Write([]byte(message))
	signature := fmt.Sprintf("%x", mac.Sum(nil))
	return signature

}
