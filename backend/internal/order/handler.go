package order

import (
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
		switch err {
		case ErrInvalidProductID:
			c.JSON(http.StatusBadRequest, gin.H{"error": utils.ErrorResponse{
				Code:    "invalid_product_id",
				Message: "invalid product id",
			}})
		case ErrProductNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": utils.ErrorResponse{
				Code:    "product_not_found",
				Message: "product not found",
			}})
		case ErrInvalidQuantity:
			c.JSON(http.StatusBadRequest, gin.H{"error": utils.ErrorResponse{
				Code:    "invalid_quantity",
				Message: "invalid quantity",
			}})
		case ErrNoItems:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": utils.ErrorResponse{
					Code:    "no_items",
					Message: "no items in order",
				},
			})
		case ErrInsufficientQuanity:
			c.JSON(http.StatusBadRequest, gin.H{"error": utils.ErrorResponse{
				Code:    "insufficient_stock",
				Message: "not enough stock available",
			}})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": utils.ErrorResponse{
				Code:    "internal_error",
				Message: "unexpected behaviour",
			}})
		}
		return
	}

	var items = make([]OrderItemResponse, 0, len(createdOrder.Items))
	for _, item := range createdOrder.Items {
		items = append(items, OrderItemResponse{
			ProductID:     item.ProductID,
			Quantity:      item.Quantity,
			PurchasePrice: item.PurchasePrice,
		})
	}

	c.JSON(http.StatusCreated, gin.H{"data": OrderResponse{
		ID:         createdOrder.ID,
		Status:     OrderStatus(createdOrder.Status),
		Items:      items,
		TotalPrice: createdOrder.TotalPrice,
		CreatedAt:  createdOrder.CreatedAt,
	}})
}

func (h *Hanlder) handleCancelOrder(c *gin.Context) {
	params := struct {
		OrderID string `uri:"id" binding:"required,uuid"`
	}{}

	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.ErrorResponse{
			Code:    "invalid_request",
			Message: "invalid request",
		}})
		return
	}

	err := h.service.Cancel(c.Request.Context(), params.OrderID)
	switch err {
	case ErrOrderNotPending:
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.ErrorResponse{
			Code:    "order_not_pending",
			Message: "order is not pending",
		}})
	case ErrOrderNotFound:
		c.JSON(http.StatusNotFound, gin.H{"error": utils.ErrorResponse{
			Code:    "order_not_found",
			Message: "order not found",
		}})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": utils.ErrorResponse{
			Code:    "internal_error",
			Message: "unexpected behaviour",
		}})
	}
	c.JSON(http.StatusOK, gin.H{"message": "order cancelled"})

}

func (h *Hanlder) handleConfirmOrder(c *gin.Context) {
	params := struct {
		OrderID string `uri:"id" binding:"required,uuid"`
	}{}

	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.ErrorResponse{
			Code:    "invalid_request",
			Message: "invalid request",
		}})
		return
	}

	err := h.service.Confirm(c.Request.Context(), params.OrderID)
	switch err {
	case ErrOrderNotPending:
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.ErrorResponse{
			Code:    "order_not_pending",
			Message: "order is not pending",
		}})
	case ErrOrderNotFound:
		c.JSON(http.StatusNotFound, gin.H{"error": utils.ErrorResponse{
			Code:    "order_not_found",
			Message: "order not found",
		}})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": utils.ErrorResponse{
			Code:    "internal_error",
			Message: "unexpected behaviour",
		}})
	}

	c.JSON(http.StatusOK, gin.H{"message": "order confirmed"})
}

func (h *Hanlder) handleGetOrder(c *gin.Context) {
	params := struct {
		OrderID string `uri:"id" binding:"required,uuid"`
	}{}

	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.ErrorResponse{
			Code:    "invalid_request",
			Message: "invalid request",
		}})
		return
	}

	order, err := h.service.GetByID(c.Request.Context(), params.OrderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": utils.ErrorResponse{
			Code:    "internal_error",
			Message: "unexpected behaviour",
		}})
		return
	}

	var items []OrderItemResponse
	for _, item := range order.Items {
		items = append(items, OrderItemResponse{
			ProductID:     item.ProductID,
			Quantity:      item.Quantity,
			PurchasePrice: item.PurchasePrice,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": OrderResponse{
		ID:         order.ID,
		Status:     OrderStatus(order.Status),
		TotalPrice: order.TotalPrice,
		Items:      items,
		CreatedAt:  order.CreatedAt,
		UpdatedAt:  order.UpdatedAt,
	}})
}

func (h *Hanlder) handleGetOrders(c *gin.Context) {
	var query utils.PaginationQuery
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
	router.POST("/orders/:id/cancel", h.handleCancelOrder)
	router.POST("/orders/:id/confirm", h.handleConfirmOrder)
	router.GET("/orders/:id", h.handleGetOrder)
	router.GET("/orders", h.handleGetOrders)
}
