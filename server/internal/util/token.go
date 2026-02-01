package util

import (
	"crypto/rand"
	"fmt"
)

const charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// GenerateToken generates a cryptographically secure random token.
// The token uses alphanumeric characters (0-9A-Za-z).
func GenerateToken(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("token length must be positive")
	}

	randomBytes := make([]byte, length)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", fmt.Errorf("failed to generate random token: %w", err)
	}

	result := make([]byte, length)
	for i := range length {
		result[i] = charset[int(randomBytes[i])%len(charset)]
	}

	return string(result), nil
}
