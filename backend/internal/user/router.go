package user

import (
	"booky-backend/internal/middleware"
	"booky-backend/pkg/config"

	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	GetUserByID(c *gin.Context)
	UpdateUser(c *gin.Context)
	DeleteUser(c *gin.Context)
	GetAllUsers(c *gin.Context)
	UserRegister(c *gin.Context)
	UserLogin(c *gin.Context)
	RefreshToken(c *gin.Context)
}

type Router struct {
	handler UserHandler
	config  *config.Config
}

func NewRouter(handler UserHandler, config *config.Config) *Router {
	return &Router{
		handler: handler,
		config:  config,
	}
}

func (r *Router) MapRoutes(vgroup *gin.RouterGroup) {
	auth := vgroup.Group("/auth")
	auth.POST("/register", r.handler.UserRegister)
	auth.POST("/login", r.handler.UserLogin)
	auth.POST("/refresh", r.handler.RefreshToken)

	users := vgroup.Group("/users")
	users.Use(middleware.Auth(r.config))
	users.GET("", r.handler.GetAllUsers)
	users.GET("/:id", r.handler.GetUserByID)
	users.PUT("/:id", r.handler.UpdateUser)
	users.DELETE("/:id", r.handler.DeleteUser)
}
