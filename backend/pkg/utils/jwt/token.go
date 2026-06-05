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

type TokenDuration time.Duration

const (
	RefreshTokenTTL   TokenDuration = TokenDuration(30 * 24 * time.Hour) // 30 days
	AccessTokenTTL    TokenDuration = TokenDuration(15 * time.Minute)
	ResetPassTokenTTL TokenDuration = TokenDuration(15 * time.Minute)
)

type Claims struct {
	jwt.RegisteredClaims
	Type TokenType `json:"type,omitempty"`
}

func CreateToken(subject, secret string, duration TokenDuration, tokenType TokenType) (string, error) {
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   subject,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(duration))),
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

func HashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
