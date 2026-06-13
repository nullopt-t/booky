package user

import (
	"booky-backend/pkg/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	UserID   string `json:"user_id"`
	UserRole string `json:"user_role"`
	Type     string `json:"type"`
	jwt.RegisteredClaims
}

type JwtService struct {
	secrets *config.Secrets
}

func NewJwtService(
	secrets *config.Secrets,
) *JwtService {
	return &JwtService{
		secrets: secrets,
	}
}

const (
	AccessTokenExpiration  = 15 * 60
	RefreshTokenExpiration = 24 * 60 * 60
)

func (s *JwtService) GenerateAccessToken(userID, userRole string) (string, error) {
	claims := UserClaims{
		UserID:   userID,
		UserRole: userRole,
		Type:     "access_token",
	}
	claims.IssuedAt = jwt.NewNumericDate(time.Now())
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(AccessTokenExpiration))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.secrets.JwtAccessTokenSecretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (s *JwtService) GenerateRefreshToken(userID string) (string, error) {
	claims := UserClaims{
		UserID: userID,
		Type:   "refresh_token",
	}
	claims.IssuedAt = jwt.NewNumericDate(time.Now())
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(RefreshTokenExpiration))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.secrets.JwtRefreshTokenSecretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (s *JwtService) GenerateTokenPair(userID, userRole string) (string, string, error) {
	accessToken, err := s.GenerateAccessToken(userID, userRole)
	if err != nil {
		return "", "", err
	}
	refreshToken, err := s.GenerateRefreshToken(userID)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}
