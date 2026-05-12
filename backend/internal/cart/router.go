package cart

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(r *gin.RouterGroup, db *pgxpool.Pool) {
	repo := NewPostgresRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	r.GET("/", handler.GetCart)
	r.POST("/items", handler.AddItem)
	r.DELETE("/", handler.EmptyCart)
}
