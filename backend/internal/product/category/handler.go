package category

import (
	"booky-backend/internal/model"
	"booky-backend/pkg/api"
	"booky-backend/pkg/database"
	"booky-backend/pkg/logger"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CategoryService interface {
	Create(ctx context.Context, name string) (*model.ProductCategory, error)
	GetAll(ctx context.Context, q *api.PageQuery) ([]*model.ProductCategory, *api.Page, error)
}

type Handler struct {
	service CategoryService
}

func NewHandler(service CategoryService) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) CreateCategory(c *gin.Context) {
	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, api.Error("BAD_REQUEST", err.Error()))
		return
	}
	category, err := h.service.Create(c.Request.Context(), req.Name)
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
	c.JSON(http.StatusCreated, CategoryResponse{
		ID:        category.ID,
		Name:      category.Name,
		CreatedAt: category.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: category.UpdatedAt.Format("2006-01-02 15:04:05"),
	})
}

func (h *Handler) GetAllCategories(c *gin.Context) {
	var req api.PageQuery
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, api.Error("BAD_REQUEST", err.Error()))
		return
	}
	categories, page, err := h.service.GetAll(c.Request.Context(), &req)
	if err != nil {
		switch err {
		case context.Canceled:
			c.JSON(http.StatusRequestTimeout, api.Error("CANCELED", err.Error()))
		case database.ErrNotFound:
			c.JSON(http.StatusRequestTimeout, api.Error("NOT_FOUND", err.Error()))
		case database.ErrConflict:
			c.JSON(http.StatusRequestTimeout, api.Error("CONFLICT", err.Error()))
		default:
			c.JSON(http.StatusInternalServerError, api.Error("INTERNAL_ERROR", err.Error()))
		}
		return
	}

	var categoriesResponse []CategoryResponse
	for _, category := range categories {
		categoriesResponse = append(categoriesResponse, CategoryResponse{
			ID:        category.ID,
			Name:      category.Name,
			CreatedAt: category.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: category.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	c.JSON(http.StatusOK, CategoriesResponse{
		Categories: categoriesResponse,
		Page:       page.Index,
		PageSize:   page.Limit,
		Total:      page.Total,
	})
}
