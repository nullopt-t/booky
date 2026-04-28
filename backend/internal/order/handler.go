package order

import (
	"booky-backend/internal/db"
	"booky-backend/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Hanlder struct {
	service *Service
}

func NewHandler(s *Service) *Hanlder {
	return &Hanlder{service: s}
}

func (h *Hanlder) handleCreateOrder(c *gin.Context) {
	var order CreateOrderRequest
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.ErrorResponse{
			Code:    "invalid_request",
			Message: err.Error(),
		}})
		return
	}

	createdOrder, err := h.service.Create(c.Request.Context(), order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": utils.ErrorResponse{
			Code:    "internal_error",
			Message: "unexpected behaviour",
		}})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": createdOrder})
}

func (h *Hanlder) handleGetOrders(c *gin.Context) {
	var query db.PaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.ErrorResponse{
			Code:    "invalid_request",
			Message: err.Error(),
		}})
		return
	}

	if query.Limit == 0 {
		query.Limit = 10
	}

	if query.Page == 0 {
		query.Page = 1
	}

	orders, err := h.service.GetAll(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": utils.ErrorResponse{
			Code:    "internal_error",
			Message: "unexpected behaviour",
		}})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": orders,
		"meta": gin.H{
			"page":  query.Page,
			"limit": query.Limit,
			"total": len(orders),
		},
	})
}

func (h *Hanlder) RegisterRoutes(router *gin.Engine) {
	router.POST("/orders", h.handleCreateOrder)
	router.GET("/orders", h.handleGetOrders)
}
