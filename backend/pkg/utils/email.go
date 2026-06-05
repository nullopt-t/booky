package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
)

func GenerateOTP() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%06d", n.Int64()), nil
}

func IsValidPhone(phone string) bool {
	return len(phone) >= 10 && len(phone) <= 15
}

func IsValidEmail(email string) bool {
	return strings.Contains(email, "@")
}
