package storage

import (
	"crypto/rand"
	mrand "math/rand"

	"github.com/tv42/zbase32"
)

const (
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func init() {
	// seed the non-crypto random number generator
	//mrand.Seed(time.Now().UnixNano())
}

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

// returns n characters encoded in zbase32 using secure random bits
func generateRandomStringForKey(length int) (string, error) {
	b, err := generateRandomBytes(length * 2)
	return zbase32.EncodeToString(b)[0:length], err
}

// GenerateRandomString generates a random string with specified length
// does NOT use hardware based crypto/rand package (should be used for IDs)
func generateRandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = letterBytes[mrand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

// generates IDs
func generateID() string {
	return generateRandomString(16)
}
