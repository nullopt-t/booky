package product

import (
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
	products, err := h.service.List(c.Request.Context())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{
				"error": utils.ErrorResponse{
					Code: "internal_error",
					Message: "unexpected behaviour",
				},
			})
		return
	}
	c.JSON(200, gin.H{
		"data": products,
	})
}

func (h *Handler) handleGetProductByID(c *gin.Context) {
	id := c.Param("id")
	p, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
	c.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{
				"error": utils.ErrorResponse{
					Code: "internal_error",
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
	c.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{
				"error": utils.ErrorResponse{
					Code: "internal_error",
					Message: "unexpected behaviour",
				},
			})
		return
	}

	newProduct, err := h.service.Create(c.Request.Context(), p)
	if err != nil {
	c.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{
				"error": utils.ErrorResponse{
					Code: "internal_error",
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
	router.GET("/products", h.handleGetProducts) router.GET("/products/:id", h.handleGetProductByID)
	router.POST("/products", h.handlerCreateProduct)
}
