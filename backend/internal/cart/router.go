package cart

import (
	"github.com/gin-gonic/gin"
)

func MapRoutes(r *gin.RouterGroup, handler CartHandler) {
	rg := r.Group("/carts")
	rg.GET("", handler.GetCart)
	rg.POST("/items", handler.AddItem)
	rg.DELETE("", handler.EmptyCart)
}
