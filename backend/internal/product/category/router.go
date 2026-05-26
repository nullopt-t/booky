package category

import (
	"github.com/gin-gonic/gin"
)

type CategoryHanlder interface {
	CreateCategory(c *gin.Context)
	GetAllCategories(c *gin.Context)
}

type Router struct {
	handler CategoryHanlder
}

func NewRouter(handler CategoryHanlder) *Router {
	return &Router{
		handler: handler,
	}
}

func (r *Router) MapRoutes(group *gin.RouterGroup) {
	group.POST("", r.handler.CreateCategory)
	group.GET("", r.handler.GetAllCategories)
}
