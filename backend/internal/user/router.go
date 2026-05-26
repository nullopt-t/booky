package user

import "github.com/gin-gonic/gin"

type UserHandler interface {
	CreateUser(c *gin.Context)
	GetUserByID(c *gin.Context)
	UpdateUser(c *gin.Context)
	DeleteUser(c *gin.Context)
	GetAllUsers(c *gin.Context)
}

type Router struct {
	handler UserHandler
}

func NewRouter(handler UserHandler) *Router {
	return &Router{
		handler: handler,
	}
}

func (r *Router) MapRoutes(group *gin.RouterGroup) {
	group.POST("", r.handler.CreateUser)
	group.GET("", r.handler.GetAllUsers)
	group.GET("/:id", r.handler.GetUserByID)
	group.PUT("/:id", r.handler.UpdateUser)
	group.DELETE("/:id", r.handler.DeleteUser)
}
