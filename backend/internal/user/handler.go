package user

import (
	"booky-backend/pkg/api"
	"booky-backend/pkg/database"
	"booky-backend/pkg/logger"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service UserService
}

func NewHandler(service UserService) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) CreateUser(c *gin.Context) {
	var user CreateUserRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, api.Error("BAD_REQUEST", err.Error()))
		return
	}
	createdUser, err := h.service.CreateUser(c.Request.Context(), user)
	if err != nil {
		logger.Log(logger.ERROR, "get cart", logger.LMeta{"error": err})
		switch err {
		case context.Canceled:
			c.JSON(http.StatusRequestTimeout, api.Error("CANCELED", err.Error()))
		case database.ErrNotFound:
			c.JSON(http.StatusRequestTimeout, api.Error("NOT_FOUND", err.Error()))
		case database.ErrConflict:
			c.JSON(http.StatusRequestTimeout, api.Error("CONFLICT", err.Error()))
		default:
			c.JSON(http.StatusRequestTimeout, api.Error("INTERNAL_ERROR", err.Error()))
		}
		return
	}
	c.JSON(http.StatusOK, api.Success(CreateUserResponse{
		ID:         createdUser.ID,
		Email:      createdUser.Email,
		IsInactive: createdUser.IsInactive,
		CreatedAt:  createdUser.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:  createdUser.UpdatedAt.Format("2006-01-02 15:04:05"),
	}))
}

func (h *Handler) GetUserByID(c *gin.Context) {
	var uri UserURIParam
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, api.Error("BAD_REQUEST", err.Error()))
		return
	}
	user, err := h.service.GetUserByID(c.Request.Context(), uri.ID)
	if err != nil {
		logger.Log(logger.ERROR, "get cart", logger.LMeta{"error": err})
		switch err {
		case context.Canceled:
			c.JSON(http.StatusRequestTimeout, api.Error("CANCELED", err.Error()))
		case database.ErrNotFound:
			c.JSON(http.StatusRequestTimeout, api.Error("NOT_FOUND", err.Error()))
		case database.ErrConflict:
			c.JSON(http.StatusRequestTimeout, api.Error("CONFLICT", err.Error()))
		default:
			c.JSON(http.StatusRequestTimeout, api.Error("INTERNAL_ERROR", err.Error()))
		}
		return
	}
	c.JSON(http.StatusOK, api.Success(UserResponse{
		ID:         user.ID,
		Email:      user.Email,
		IsInactive: user.IsInactive,
		CreatedAt:  user.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:  user.UpdatedAt.Format("2006-01-02 15:04:05"),
	}))
}

func (h *Handler) UpdateUser(c *gin.Context) {
	var uri UserURIParam
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, api.Error("BAD_REQUEST", err.Error()))
		return
	}
	var user UpdateUserRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, api.Error("BAD_REQUEST", err.Error()))
		return
	}

	err := h.service.UpdateUser(c.Request.Context(), uri.ID, &user)
	if err != nil {
		logger.Log(logger.ERROR, "get cart", logger.LMeta{"error": err})
		switch err {
		case context.Canceled:
			c.JSON(http.StatusRequestTimeout, api.Error("CANCELED", err.Error()))
		case database.ErrNotFound:
			c.JSON(http.StatusRequestTimeout, api.Error("NOT_FOUND", err.Error()))
		case database.ErrConflict:
			c.JSON(http.StatusRequestTimeout, api.Error("CONFLICT", err.Error()))
		default:
			c.JSON(http.StatusRequestTimeout, api.Error("INTERNAL_ERROR", err.Error()))
		}
		return
	}
	c.JSON(http.StatusOK, api.Success("User updated successfully"))
}

func (h *Handler) GetAllUsers(c *gin.Context) {
	var q api.PageQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		c.JSON(http.StatusBadRequest, api.Error("BAD_REQUEST", err.Error()))
		return
	}
	users, page, err := h.service.GetAllUsers(c.Request.Context(), &q)
	if err != nil {
		logger.Log(logger.ERROR, "get cart", logger.LMeta{"error": err})
		switch err {
		case context.Canceled:
			c.JSON(http.StatusRequestTimeout, api.Error("CANCELED", err.Error()))
		case database.ErrNotFound:
			c.JSON(http.StatusRequestTimeout, api.Error("NOT_FOUND", err.Error()))
		case database.ErrConflict:
			c.JSON(http.StatusRequestTimeout, api.Error("CONFLICT", err.Error()))
		default:
			c.JSON(http.StatusRequestTimeout, api.Error("INTERNAL_ERROR", err.Error()))
		}
		return
	}

	data := make([]UserResponse, 0, len(users))
	for _, user := range users {
		data = append(data, UserResponse{
			ID:         user.ID,
			Email:      user.Email,
			IsInactive: user.IsInactive,
			CreatedAt:  user.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:  user.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	c.JSON(http.StatusOK, api.SuccessPaginated(data, page))
}

func (h *Handler) DeleteUser(c *gin.Context) {
	var uri UserURIParam
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, api.Error("BAD_REQUEST", err.Error()))
		return
	}
	err := h.service.DeleteUser(c.Request.Context(), uri.ID)
	if err != nil {
		logger.Log(logger.ERROR, "get cart", logger.LMeta{"error": err})
		switch err {
		case context.Canceled:
			c.JSON(http.StatusRequestTimeout, api.Error("CANCELED", err.Error()))
		case database.ErrNotFound:
			c.JSON(http.StatusRequestTimeout, api.Error("NOT_FOUND", err.Error()))
		case database.ErrConflict:
			c.JSON(http.StatusRequestTimeout, api.Error("CONFLICT", err.Error()))
		default:
			c.JSON(http.StatusRequestTimeout, api.Error("INTERNAL_ERROR", err.Error()))
		}
		return
	}
	c.JSON(http.StatusOK, api.Success("User deleted successfully"))
}
