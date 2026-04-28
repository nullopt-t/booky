package product

import (
	"booky-backend/internal/db"
	"booky-backend/internal/utils"
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
	var query db.PaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest,
			gin.H{
				"error": utils.ErrorResponse{
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

	ppr, err := h.service.GetAll(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{
				"error": utils.ErrorResponse{
					Code:    "internal_error",
					Message: "unexpected behaviour",
				},
			})
		return
	}

	c.JSON(200, gin.H{
		"data": ppr.Products,
		"meta": gin.H{
			"page":  ppr.Page,
			"limit": ppr.Limit,
			"total": ppr.Total,
		},
	})
}

func (h *Handler) handleGetProductByID(c *gin.Context) {
	var params = struct {
		ID string `uri:"id" binding:"required"`
	}{}

	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest,
			gin.H{
				"error": utils.ErrorResponse{
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
				"error": utils.ErrorResponse{
					Code:    "internal_error",
					Message: "unexpected behaviour",
				},
			})
		return
	}

	c.JSON(http.StatusFound, gin.H{
		"data": p,
	})
}

func (h *Handler) handlerCreateProduct(c *gin.Context) {
	var p CreateProductRequest
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{
				"error": utils.ErrorResponse{
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
				"error": utils.ErrorResponse{
					Code:    "internal_error",
					Message: "unexpected behaviour",
				},
			})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": newProduct,
	})
}

func (h *Handler) RegisterRoutes(router *gin.Engine) {
	router.GET("/products", h.handleGetProducts)
	router.GET("/products/:id", h.handleGetProductByID)
	router.POST("/products", h.handlerCreateProduct)
}
