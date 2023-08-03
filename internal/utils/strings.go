package utils

import (
	"math/rand"
	"strings"
)

func RandString(n int, randomChars string) string {
	runes := []rune(randomChars)
	b := make([]rune, n)
	for i := range b {
		b[i] = runes[rand.Intn(len(runes))]
	}
	return string(b)
}

// Contains checks if a string is present in a slice
func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

// AddNewLine changes a string by adding a newline character ("\n") every n characters
func AddNewLine(str string, n int) string {
	var sb strings.Builder

	for i, char := range str {
		sb.WriteRune(char)

		if (i+1)%n == 0 {
			sb.WriteString("\n")
		}
	}

	return sb.String()
}
