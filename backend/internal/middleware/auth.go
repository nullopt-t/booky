package middleware

import (
	"booky-backend/pkg/api"
	"booky-backend/pkg/config"
	"booky-backend/pkg/utils/jwt"
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

func Auth(config *config.Config) gin.HandlerFunc {
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

		c.Set("userID", claims.Subject)
		c.Set("claims", claims)
		c.Next()
	}
}
