package cart

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, handler CartHandler) {
	r.GET("/", handler.GetCart)
	r.POST("/items", handler.AddItem)
	r.DELETE("/", handler.EmptyCart)
}
