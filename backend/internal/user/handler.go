package user

import (
	"booky-backend/pkg/api"
	"booky-backend/pkg/config"
	"booky-backend/pkg/database"
	"booky-backend/pkg/logger"
	"booky-backend/pkg/utils"
	"booky-backend/pkg/utils/jwt"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service UserService
	config  *config.Config
}

func NewHandler(service UserService, config *config.Config) *Handler {
	return &Handler{
		service: service,
		config:  config,
	}
}

func (h *Handler) handlerError(c *gin.Context, err error) {
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
}

func (h *Handler) UserRegister(c *gin.Context) {
	var credentials RegisterUserRequest
	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, api.Error("BAD_REQUEST", err.Error()))
		return
	}

	user, err := h.service.CreateUser(c.Request.Context(), CreateUserRequest(credentials))
	if err != nil {
		h.handlerError(c, err)
		return
	}

	token, err := jwt.CreateToken(user.ID.String(), h.config.JwtSecretKey, jwt.AccessTokenTTL, jwt.Access)
	if err != nil {
		c.JSON(http.StatusInternalServerError, api.Error("INTERNAL_ERROR", "interval error"))
	}

	c.JSON(http.StatusOK, api.Success(RegisterUserResponse{
		Email:       user.Email,
		AccessToken: token,
	}))
}

func (h *Handler) UserLogin(c *gin.Context) {
	var credentials UserCredentials
	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, api.Error("BAD_REQUEST", err.Error()))
		return
	}

	user, err := h.service.GetUserByEmail(c.Request.Context(), credentials.Email)
	if err != nil {
		h.handlerError(c, err)
		return
	}

	if err := utils.ComparePassword(user.PasswordHash, credentials.Password); err != nil {
		h.handlerError(c, err)
		return
	}

	accessToken, err := jwt.CreateToken(user.ID.String(), h.config.JwtSecretKey, jwt.AccessTokenTTL, jwt.Access)
	if err != nil {
		logger.Log(logger.ERROR, err.Error())
		c.JSON(http.StatusInternalServerError, api.Error("INTERNAL_ERROR", "interval error"))
		return
	}

	refreshToken, err := jwt.CreateToken(user.ID.String(), h.config.JwtSecretKey, jwt.RefreshTokenTTL, jwt.Refresh)
	if err != nil {
		logger.Log(logger.ERROR, err.Error())
		c.JSON(http.StatusInternalServerError, api.Error("INTERNAL_ERROR", "interval error"))
		return
	}

	c.JSON(http.StatusOK, api.Success(RegisterUserResponse{
		Email:        user.Email,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
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
		h.handlerError(c, err)
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
		h.handlerError(c, err)
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

	if q.Limit == 0 {
		q.Limit = 5
	}

	if q.Page == 0 {
		q.Page = 1
	}

	users, page, err := h.service.GetAllUsers(c.Request.Context(), &q)
	if err != nil {
		h.handlerError(c, err)
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
		h.handlerError(c, err)
		return
	}
	c.JSON(http.StatusOK, api.Success("User deleted successfully"))
}

func (h *Handler) RefreshToken(c *gin.Context) {
	var body RefreshTokenRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, api.Error("BAD_REQUEST", err.Error()))
		return
	}

	claims, err := jwt.VerifyToken(body.Refresh_token, h.config.JwtSecretKey)
	if err != nil {
		c.JSON(http.StatusForbidden, api.Error("INVALID_REFRESH_TOKEN", err.Error()))
		return
	}

	accessToken, err := jwt.CreateToken(claims.Subject, h.config.JwtSecretKey, jwt.AccessTokenTTL, jwt.Access)
	if err != nil {
		c.JSON(http.StatusInternalServerError, api.Error("INTERNAL_ERROR", err.Error()))
		return
	}

	c.JSON(http.StatusOK, api.Success(RefreshTokenResponse{
		AccessToken: accessToken,
	}))
}
