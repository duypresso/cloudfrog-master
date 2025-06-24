package utils

import (
	"crypto/rand"
	"math/big"
	"strings"
)

const (
	base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	codeLength  = 6
)

// GenerateShortCode creates a random Base62 shortcode
func GenerateShortCode() (string, error) {
	var result strings.Builder

	for i := 0; i < codeLength; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(base62Chars))))
		if err != nil {
			return "", err
		}
		result.WriteByte(base62Chars[n.Int64()])
	}

	return result.String(), nil
}
