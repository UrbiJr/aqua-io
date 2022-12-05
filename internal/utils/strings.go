package utils

import "math/rand"

func RandString(n int, randomChars string) string {
	runes := []rune(randomChars)
	b := make([]rune, n)
	for i := range b {
		b[i] = runes[rand.Intn(len(runes))]
	}
	return string(b)
}
