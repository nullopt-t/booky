package checkout

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(handler CheckoutHandler, r *gin.RouterGroup, db *pgxpool.Pool) {
	r.POST("/", handler.HandleCheckout)
}
