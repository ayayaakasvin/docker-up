package password

import (
	"testing"
)

type test struct {
	password string
	expected bool
}

func TestIsValidPassword(t *testing.T) {
	tests := []test{
		{
			password: "",
			expected: false,
		},
		{
			password: "short",
			expected: false,
		},
		{
			password: "alllowercase",
			expected: false,
		},
		{
			password: "ALLUPPERCASE",
			expected: false,
		},
		{
			password: "12345678",
			expected: false,
		},
		{
			password: "Valid123",
			expected: true,
		},
		{
			password: "Another1Valid",
			expected: true,
		},
		{
			password: "NoDigitsHere!",
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.password, func(t *testing.T) {
			result := IsValidPassword(tc.password)
			if result != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}