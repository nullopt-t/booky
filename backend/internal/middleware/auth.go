package middleware

import (
	"booky-backend/internal/model"
	"booky-backend/internal/shared/token"
	"booky-backend/pkg/api"
	"booky-backend/pkg/config"
	"booky-backend/pkg/utils/jwt"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const prefix = "Bearer "

func GetClaims(c *gin.Context) *jwt.Claims {
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
		value, ok := c.Get("user")
		if !ok {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				api.Error(
					"MISSING_USER_ID",
					"missing user id",
				),
			)
			return
		}

		userSubject, ok := value.(token.UserSubject)
		if !ok {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				api.Error(
					"MISSING_USER_ID",
					"missing user id",
				),
			)
			return
		}

		if userSubject.UserRole != requiredRole {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				api.Error(
					"FORBIDDEN",
					"no permission",
				),
			)
			return
		}

		c.Next()
	}
}

func Authanticate(config *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				api.Error(
					"MISSING_AUTH_HEADER",
					"authorization header is required",
				),
			)
			return
		}

		if !strings.HasPrefix(authHeader, prefix) {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				api.Error(
					"INVALID_AUTH_FORMAT",
					"authorization header must use Bearer scheme",
				),
			)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, prefix)

		claims, err := jwt.VerifyToken(
			tokenString,
			config.JwtSecretKey,
		)

		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				api.Error(
					"INVALID_TOKEN",
					err.Error(),
				),
			)
			return
		}

		var userSubject token.UserSubject
		if err := json.Unmarshal([]byte(claims.Subject), &userSubject); err != nil {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				api.Error(
					"INVALID_TOKEN",
					err.Error(),
				),
			)
			return
		}

		c.Set("user", userSubject)
		c.Set("claims", claims)
		c.Next()
	}
}
