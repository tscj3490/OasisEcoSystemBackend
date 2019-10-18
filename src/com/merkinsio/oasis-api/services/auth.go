package services

import (
	"math/rand"
	"time"
)

var defaultCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
var seededRand = rand.New(rand.NewSource(time.Now().UTC().UnixNano()))

/*GeneratePasswordWithCharset Generate random password with the given length and charset*/
func GeneratePasswordWithCharset(length int, charset string) string {
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[seededRand.Intn(len(charset))]
	}

	return string(result)
}

/*GeneratePassword Generate random passworth with the given length using the default alphanumeric charset*/
func GeneratePassword(length int) string {
	return GeneratePasswordWithCharset(length, defaultCharset)
}
