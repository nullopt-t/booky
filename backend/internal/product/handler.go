package product

import (
	"booky-backend/internal/trans"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

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
// @Param query query trans.PaginationQuery true "Pagination query"
// @Success 200 {object} ProductsResponse
// @Failure 400 {object} trans.ApiErr
// @Failure 500 {object} trans.ApiErr
// @Router /products [get]
func (h *Handler) GetAllProducts(c *gin.Context) {
	var query trans.PaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest,
			gin.H{
				"error": trans.ApiErr{
					Code:    "invalid_request",
					Message: err.Error(),
				},
			})
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
			gin.H{
				"error": trans.ApiErr{
					Code:    "internal_error",
					Message: "unexpected behaviour",
				},
			})
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
			gin.H{
				"error": trans.ApiErr{
					Code:    "invalid_request",
					Message: err.Error(),
				},
			})
		return
	}

	productID, err := uuid.Parse(params.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			gin.H{
				"error": trans.ApiErr{
					Code:    "invalid_request",
					Message: err.Error(),
				},
			})
		return
	}

	p, err := h.service.GetByID(c.Request.Context(), productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{
				"error": trans.ApiErr{
					Code:    "internal_error",
					Message: "unexpected behaviour",
				},
			})
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
			gin.H{
				"error": trans.ApiErr{
					Code:    "invalid_request",
					Message: err.Error(),
				},
			})
		return
	}

	newProduct, err := h.service.Create(c.Request.Context(), p)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{
				"error": trans.ApiErr{
					Code:    "internal_error",
					Message: "unexpected behaviour",
				},
			})
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
