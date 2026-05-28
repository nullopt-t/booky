package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenType string

const (
	Access  TokenType = "access"
	Refresh TokenType = "refresh"
)

const RefreshTokenTTL = 30 * 24 * time.Hour // Month
const AccessTokenTTL = 15 * time.Minute

type Claims struct {
	jwt.RegisteredClaims
	Type TokenType `json:"type,omitempty"`
}

func CreateToken(userID, secret string, expireAfter time.Duration, tokenType TokenType) (string, error) {
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expireAfter)),
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
