package order

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(router *gin.RouterGroup, db *pgxpool.Pool) {
	orderRepo := NewPostgresRepo(db)
	orderService := NewService(orderRepo)
	orderHandler := NewHandler(orderService)
	router.POST("/", orderHandler.CreateOrder)
	router.POST("/:id/cancel", orderHandler.CancelOrder)
	router.POST("/:id/confirm", orderHandler.ConfirmOrder)
	router.GET("/:id", orderHandler.GetOrderByID)
	router.GET("/", orderHandler.GetAllOrders)
}
