package service

import (
	"crypto/rand"
	"errors"
	"math/big"
)

const (
	// Default code length
	DefaultCodeLength = 6
)

// GenerateShortCode generates a random base62 encoded string
func GenerateShortCode(length int, chars string) (string, error) {
	if length <= 0 {
		length = DefaultCodeLength
	}

	if chars == "" {
		return "", errors.New("base62 characters cannot be empty")
	}

	result := make([]byte, length)
	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			return "", err
		}
		result[i] = chars[num.Int64()]
	}

	return string(result), nil
}
