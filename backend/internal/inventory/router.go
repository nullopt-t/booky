package inventory

import "github.com/gin-gonic/gin"

type InventoryHandler interface {
	GetAvailable(c *gin.Context)
	GetReserved(c *gin.Context)
}

type Router struct {
	handler InventoryHandler
}

func NewRouter(handler InventoryHandler) *Router {
	return &Router{
		handler: handler,
	}
}

func (r *Router) MapRoutes(group *gin.RouterGroup) {
	group.GET("/:product_id/available", r.handler.GetAvailable)
	group.GET("/:product_id/reserved", r.handler.GetReserved)
}
