package product

import (
	"booky-backend/internal/trans"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) handleGetProducts(c *gin.Context) {
	var query trans.PaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest,
			gin.H{
				"error": trans.ErrorResponse{
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
				"error": trans.ErrorResponse{
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
			Stock:     p.Stock,
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

func (h *Handler) handleGetProductByID(c *gin.Context) {
	var params = struct {
		ID string `uri:"id" binding:"required"`
	}{}

	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest,
			gin.H{
				"error": trans.ErrorResponse{
					Code:    "invalid_request",
					Message: err.Error(),
				},
			})
		return
	}

	p, err := h.service.GetByID(c.Request.Context(), params.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{
				"error": trans.ErrorResponse{
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
		Stock:     p.Stock,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	})
}

func (h *Handler) handlerCreateProduct(c *gin.Context) {
	var p CreateProductRequest
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{
				"error": trans.ErrorResponse{
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
				"error": trans.ErrorResponse{
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
		Stock:     newProduct.Stock,
		CreatedAt: newProduct.CreatedAt,
		UpdatedAt: newProduct.UpdatedAt,
	})
}

func (h *Handler) handlerUpdateProduct(c *gin.Context) {
	var params = struct {
		ID string `uri:"id" binding:"required"`
	}{}

	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest,
			gin.H{
				"error": trans.ErrorResponse{
					Code:    "invalid_request",
					Message: err.Error(),
				},
			})
		return
	}

	var p UpdateProductRequest
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{
				"error": trans.ErrorResponse{
					Code:    "invalid_request",
					Message: err.Error(),
				},
			})
		return
	}

	updatedProduct, err := h.service.Update(c.Request.Context(), params.ID, p)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{
				"error": trans.ErrorResponse{
					Code:    "internal_error",
					Message: "unexpected behaviour",
				},
			})
		return
	}

	c.JSON(http.StatusOK, &ProductResponse{
		ID:        updatedProduct.ID,
		Title:     updatedProduct.Title,
		Price:     updatedProduct.Price,
		Stock:     updatedProduct.Stock,
		CreatedAt: updatedProduct.CreatedAt,
		UpdatedAt: updatedProduct.UpdatedAt,
	})
}

func (h *Handler) RegisterRoutes(router *gin.Engine) {
	router.GET("/products", h.handleGetProducts)
	router.GET("/products/:id", h.handleGetProductByID)
	router.POST("/products", h.handlerCreateProduct)
	router.PATCH("/products/:id", h.handlerUpdateProduct)
}
