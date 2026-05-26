package payment

import "github.com/gin-gonic/gin"

func MapRoutes(r *gin.RouterGroup, h *Handler) {
	rg := r.Group("/payments")
	rg.GET("", nil)
	rg.POST("", nil)
}
