package order

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, orderHandler OrderHandler) {
	router.POST("/", orderHandler.CreateOrder)
	router.POST("/:id/cancel", orderHandler.CancelOrder)
	router.POST("/:id/confirm", orderHandler.ConfirmOrder)
	router.GET("/:id", orderHandler.GetOrderByID)
	router.GET("/", orderHandler.GetAllOrders)
}
