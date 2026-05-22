package inventory

import "github.com/gin-gonic/gin"

func MapRoutes(r *gin.RouterGroup, h InventoryHandler) {
	rg := r.Group("/inventories")
	rg.GET("/:product_id/available", h.GetAvailable)
	rg.GET("/:product_id/reserved", h.GetReserved)
}
