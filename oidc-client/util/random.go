package util

import (
	"math/rand"
	"time"
)

const (
	allowedChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
)

func init() {
	// seed math/rand with current unix time with nano seconds precision
	rand.Seed(time.Now().UnixNano())
}

// GenerateRandomString generates a random string with specified length
func GenerateRandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = allowedChars[rand.Int63()%int64(len(allowedChars))]
	}
	return string(b)
}
