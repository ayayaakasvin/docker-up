package hashtool

import (
	"crypto/rand"
	"fmt"
)

// GenerateKey is a function that generates a key for the given integer value
func GenerateKey(size int) ([]byte, error) {
	key := make([]byte, size)
	_, err := rand.Read(key)
	if err != nil {
		return nil, fmt.Errorf("failed to generate key: %w", err)
	}

	return key, nil
}