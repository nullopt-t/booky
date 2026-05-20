package product

import "github.com/gin-gonic/gin"

func MapRoutes(r *gin.RouterGroup, h ProductHandler) {
	rg := r.Group("/products")
	rg.GET("/", h.GetAllProducts)
	rg.GET("/:id", h.GetProductByID)
	rg.POST("/", h.CreateProduct)
}
