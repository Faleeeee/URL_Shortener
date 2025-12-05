package service

import (
	"crypto/rand"
	"math/big"
)

const (
	// Base62 characters for encoding
	base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	// Default code length
	DefaultCodeLength = 6
)

// GenerateShortCode generates a random base62 encoded string
func GenerateShortCode(length int) (string, error) {
	if length <= 0 {
		length = DefaultCodeLength
	}

	result := make([]byte, length)
	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(base62Chars))))
		if err != nil {
			return "", err
		}
		result[i] = base62Chars[num.Int64()]
	}

	return string(result), nil
}
