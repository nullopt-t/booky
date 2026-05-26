package order

import (
	"github.com/gin-gonic/gin"
)

type OrderHandler interface {
	GetOrderByID(c *gin.Context)
	GetAllOrders(c *gin.Context)
	CancelOrder(c *gin.Context)
	ConfirmOrder(c *gin.Context)
}

type Router struct {
	handler OrderHandler
}

func NewRouter(handler OrderHandler) *Router {
	return &Router{handler: handler}
}

func (r *Router) MapRoutes(group *gin.RouterGroup) {
	group.POST("/:id/cancel", r.handler.CancelOrder)
	group.POST("/:id/confirm", r.handler.ConfirmOrder)
	group.GET("/:id", r.handler.GetOrderByID)
	group.GET("", r.handler.GetAllOrders)
}
