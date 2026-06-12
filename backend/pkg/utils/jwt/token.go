package jwt

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenType string

const (
	AccessTokenType    TokenType = "access"
	RefreshTokenType   TokenType = "refresh"
	ResetPassTokenType TokenType = "reset_pass"
)

const (
	RefreshTokenTTL   = 30 * 24 * time.Hour
	AccessTokenTTL    = 15 * time.Minute
	ResetPassTokenTTL = 48 * time.Hour
)

type Claims struct {
	jwt.RegisteredClaims
	Type TokenType `json:"type,omitempty"`
}

func CreateToken(subject, secret string, duration time.Duration, tokenType TokenType) (string, error) {
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   subject,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
		Type: tokenType,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func VerifyToken(tokenStr, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (any, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("unexpected signing method %s", t.Method.Alg())
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if token == nil {
		return nil, fmt.Errorf("nil token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, fmt.Errorf("invalid claims type")
	}
	return claims, nil
}

func Hash(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
