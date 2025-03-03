package password

import (
	"regexp"
)

var (
	uppercase = regexp.MustCompile(`[A-Z]`)
	lowercase = regexp.MustCompile(`[a-z]`)
	digit     = regexp.MustCompile(`[0-9]`)
)

// IsValidPassword checks if the password is valid
func IsValidPassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	if !uppercase.MatchString(password) {
		return false
	}

	if !lowercase.MatchString(password) {
		return false
	}

	if !digit.MatchString(password) {
		return false
	}

	return true
}
