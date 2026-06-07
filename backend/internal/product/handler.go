package product

import (
	"booky-backend/internal/model"
	"booky-backend/pkg/api"
	"booky-backend/pkg/database"
	"booky-backend/pkg/logger"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProudctService interface {
	Create(ctx context.Context, req CreateProductRequest) (*model.Product, error)
	Update(ctx context.Context, productID uuid.UUID, req UpdateProductRequest) (*model.Product, error)
	GetByID(ctx context.Context, productID uuid.UUID) (*model.Product, error)
	GetAll(ctx context.Context, q api.PageQuery) ([]*model.Product, *api.Page, error)
	CreateCategory(ctx context.Context, name string) (*model.ProductCategory, error)
	GetAllCategories(ctx context.Context, q *api.PageQuery) ([]*model.ProductCategory, *api.Page, error)
}

type Handler struct {
	service ProudctService
}

func NewHandler(s ProudctService) ProductHandler {
	return &Handler{s}
}

// GetAllProducts godoc
// @Summary Get all products
// @Description Get all products
// @Tags Products
// @Accept json
// @Produce json
// @Param query query api.PageQuery true "Pagination query"
// @Success 200 {object} ProductsResponse
// @Failure 400 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /products [get]
func (h *Handler) GetAllProducts(c *gin.Context) {
	var query api.PageQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest,
			api.Error("invalid_request", err.Error()))
		return
	}

	if query.Limit == 0 {
		query.Limit = 10
	}

	if query.Page == 0 {
		query.Page = 1
	}

	result, page, err := h.service.GetAll(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			api.Error("invalid_request", err.Error()))
		return
	}

	var products []ProductResponse
	for _, p := range result {
		products = append(products, ProductResponse{
			ID:        p.ID,
			Title:     p.Title,
			Price:     p.Price,
			CreatedAt: p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
		})
	}

	c.JSON(200, &ProductsResponse{
		Products: products,
		Page:     page.Index,
		Limit:    page.Limit,
		Total:    page.Total,
	})
}

func (h *Handler) GetProductByID(c *gin.Context) {
	var params = struct {
		ID string `uri:"id" binding:"required,uuid"`
	}{}

	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest,
			api.Error("invalid_request", err.Error()))
		return
	}

	productID, err := uuid.Parse(params.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			api.Error("invalid_request", err.Error()))
		return
	}

	p, err := h.service.GetByID(c.Request.Context(), productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			api.Error("internal_error", "unexpected behaviour"))
		return
	}

	c.JSON(http.StatusFound, &ProductResponse{
		ID:        p.ID,
		Title:     p.Title,
		Price:     p.Price,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	})
}

func (h *Handler) CreateProduct(c *gin.Context) {
	var p CreateProductRequest
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusInternalServerError,
			api.Error("invalid_request", err.Error()))
		return
	}

	newProduct, err := h.service.Create(c.Request.Context(), p)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			api.Error("internal_error", "unexpected behaviour"))
		return
	}

	c.JSON(http.StatusCreated, &ProductResponse{
		ID:        newProduct.ID,
		Title:     newProduct.Title,
		Price:     newProduct.Price,
		CreatedAt: newProduct.CreatedAt,
		UpdatedAt: newProduct.UpdatedAt,
	})
}

func (h *Handler) CreateCategory(c *gin.Context) {
	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, api.Error("BAD_REQUEST", err.Error()))
		return
	}
	category, err := h.service.CreateCategory(c.Request.Context(), req.Name)
	if err != nil {
		logger.Log(logger.ERROR, "get cart", logger.Meta{"error": err})
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
	categories, page, err := h.service.GetAllCategories(c.Request.Context(), &req)
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
