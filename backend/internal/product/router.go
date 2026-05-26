package product

import "github.com/gin-gonic/gin"


type ProductHandler interface {
	CreateProduct(c *gin.Context)
	GetProductByID(c *gin.Context)
	GetAllProducts(c *gin.Context)
}


type Router struct {
	handler ProductHandler
}

func NewRouter(handler ProductHandler) *Router {
	return &Router{
		handler: handler,
	}
}

func (r *Router) MapRoutes(rg *gin.RouterGroup) {
	rg.GET("", r.handler.GetAllProducts)
	rg.GET("/:id", r.handler.GetProductByID)
	rg.POST("", r.handler.CreateProduct)
}
