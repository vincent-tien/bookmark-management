package utils

import (
	"crypto/rand"
	"errors"
)

const alphaNumeric = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// GenerateRandomString generates a random alphanumeric string of the specified length.
// It uses cryptographically secure random number generation.
// Returns the generated string and an error if the length is invalid or generation fails.
func GenerateRandomString(n int) (string, error) {
	if n <= 0 {
		return "", errors.New("invalid length")
	}

	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	result := make([]byte, n)
	for i := range b {
		result[i] = alphaNumeric[int(b[i])%len(alphaNumeric)]
	}

	return string(result), nil
}
