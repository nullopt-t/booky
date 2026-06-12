package middleware

import (
	"booky-backend/internal/model"
	"booky-backend/internal/shared/token"
	"booky-backend/pkg/api/security"
	"booky-backend/pkg/config"
	"booky-backend/pkg/log"
	"booky-backend/pkg/utils/jwt"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	ErrKeyNotFound        = errors.New("key not found in the context")
	ErrInvalidUserSubject = errors.New("invalid user subject")
)

const prefix = "Bearer "

func GetUserWithContext(c *gin.Context) (*token.UserSubject, error) {
	value, ok := c.Get("user")
	if !ok {
		return nil, ErrKeyNotFound
	}

	userSubject, ok := value.(token.UserSubject)
	if !ok {
		return nil, ErrInvalidUserSubject
	}

	return &userSubject, nil
}

func GetClaimsWithContext(c *gin.Context) *jwt.Claims {
	claims, ok := c.Get("claims")
	if !ok {
		return &jwt.Claims{}
	}
	tclaims, ok := claims.(*jwt.Claims)
	if !ok {
		return &jwt.Claims{}
	}
	return tclaims
}

func Authorize(requiredRole model.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := GetUserWithContext(c)
		if err != nil {
			c.Error(
				security.NewSecureError(
					http.StatusUnauthorized,
					"INVALID_USER",
					"invalid user",
					err,
				),
			)
			c.Abort()
			return
		}
		if user.UserRole != requiredRole {
			c.Error(
				security.NewSecureError(
					http.StatusUnauthorized,
					"FORBIDDEN",
					"no permission",
					nil,
				),
			)
			c.Abort()
			return
		}

		c.Next()
	}
}

func Authanticate(secrets *config.Secrets) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.Error(
				security.NewSecureError(
					http.StatusUnauthorized,
					"MISSING_AUTH_HEADER",
					"authorization header is required",
					nil,
				),
			)
			c.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, prefix) {
			c.Error(
				security.NewSecureError(
					http.StatusUnauthorized,
					"INVALID_AUTH_FORMAT",
					"authorization header must use Bearer scheme",
					nil,
				),
			)
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, prefix)

		claims, err := jwt.VerifyToken(
			tokenString,
			secrets.JwtAccessTokenSecretKey,
		)

		if err != nil {
			c.Error(
				security.NewSecureError(
					http.StatusUnauthorized,
					"INVALID_TOKEN",
					"invalid or expired token",
					err,
				),
			)
			c.Abort()
			return
		}

		var userSubject token.UserSubject
		if err := json.Unmarshal([]byte(claims.Subject), &userSubject); err != nil {
			c.Error(
				security.NewSecureError(
					http.StatusUnauthorized,
					"INVALID_TOKEN_SUBJECT",
					"invalid token subject",
					err,
				),
			)
			c.Abort()
			return
		}

		fmt.Println("userSubject:", userSubject)
		c.Set("user", userSubject)
		c.Set("claims", claims)
		c.Next()
	}
}

func IsEmailVerified(logger log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := GetUserWithContext(c)
		if err != nil {
			c.Error(
				security.NewSecureError(
					http.StatusUnauthorized,
					"UNAUTHORIZED",
					"user not authenticated",
					err,
				),
			)
			c.Abort()
			return
		}

		if !user.IsEmailVerified {
			c.Error(
				security.NewSecureError(
					http.StatusUnauthorized,
					"UNAUTHORIZED",
					"user not authenticated: email not verified",
					nil,
				),
			)
			c.Abort()
			return
		}

		c.Next()
	}
}
