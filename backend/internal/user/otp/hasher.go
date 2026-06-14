package otp

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashOTP(otp string) (string, error) {
	hash := sha256.New()
	_, err := hash.Write([]byte(otp))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}
