package checkout

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, handler CheckoutHandler) {
	rg := r.Group("/checkout")
	rg.POST("", handler.HandleCheckout)
}
