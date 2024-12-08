package utils

import (
	"crypto/rand"
	"fmt"
)

func GenerateConfirmationCode() string {
	// Generate a random 6-digit number
	code := make([]byte, 3) // 3 bytes for 6 hex digits
	_, err := rand.Read(code)
	if err != nil {
		panic(err) // Handle error appropriately
	}
	return fmt.Sprintf("%06x", code)
}
