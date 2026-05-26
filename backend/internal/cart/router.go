package cart

import (
	"github.com/gin-gonic/gin"
)

type CartHandler interface {
	GetCart(c *gin.Context)
	AddItem(c *gin.Context)
	EmptyCart(c *gin.Context)
}

type Router struct {
	handler CartHandler
}

func NewRouter(handler CartHandler) *Router {
	return &Router{
		handler: handler,
	}
}

func (r *Router) MapRoutes(group *gin.RouterGroup) {
	group.GET("", r.handler.GetCart)
	group.POST("/items", r.handler.AddItem)
	group.DELETE("", r.handler.EmptyCart)
}
