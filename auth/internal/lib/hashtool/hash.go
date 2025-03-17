package hashtool

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// BcryptHashing is a function that hashes a password using bcrypt
func BcryptHashing(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return string(hashed), nil
}

// BcryptCompare is a function that checks if the provided password matches the hashed password
func BcryptCompare(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
