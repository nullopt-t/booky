package order

import (
	"github.com/gin-gonic/gin"
)

func MapRoutes(r *gin.RouterGroup, orderHandler OrderHandler) {
	rg := r.Group("/orders")
	rg.POST("/:id/cancel", orderHandler.CancelOrder)
	rg.POST("/:id/confirm", orderHandler.ConfirmOrder)
	rg.GET("/:id", orderHandler.GetOrderByID)
	rg.GET("", orderHandler.GetAllOrders)
}
