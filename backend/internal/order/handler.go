package order

import (
	"booky-backend/internal/model"
	"booky-backend/pkg/api"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type OrderService interface {
	GetByID(ctx context.Context, orderID uuid.UUID) (*model.Order, error)
	GetAll(ctx context.Context, q *api.PageQuery) ([]*model.Order, *api.Page, error)
	Cancel(ctx context.Context, orderID uuid.UUID) error
	Confirm(ctx context.Context, orderID uuid.UUID) error
}

type Hanlder struct {
	service OrderService
}

func NewHandler(s OrderService) OrderHandler {
	return &Hanlder{service: s}
}

func (h *Hanlder) CancelOrder(c *gin.Context) {
	var params = struct {
		ID string `uri:"id" binding:"required,uuid"`
	}{}

	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest,
			api.Error("invalid_request", err.Error()))
		return
	}

	orderID, err := uuid.Parse(params.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			api.Error("invalid_request", err.Error()))
		return
	}

	err = h.service.Cancel(c.Request.Context(), orderID)
	switch err {
	case ErrOrderNotPending:
		c.JSON(http.StatusBadRequest, api.Error("order_not_pending", "order is not pending"))
	case ErrOrderNotFound:
		c.JSON(http.StatusNotFound, api.Error("order_not_found", "order not found"))
	default:
		c.JSON(http.StatusInternalServerError, api.Error("internal_error", "unexpected behaviour"))
	}
	c.JSON(http.StatusOK, gin.H{"message": "order cancelled"})

}

func (h *Hanlder) ConfirmOrder(c *gin.Context) {
	var params = struct {
		ID string `uri:"id" binding:"required,uuid"`
	}{}

	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest,
			api.Error("invalid_request", err.Error()))
		return
	}

	orderID, err := uuid.Parse(params.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			api.Error("invalid_request", err.Error()))
		return
	}

	err = h.service.Confirm(c.Request.Context(), orderID)
	switch err {
	case ErrOrderNotPending:
		c.JSON(http.StatusBadRequest, api.Error("order_not_pending", "order is not pending"))
	case ErrOrderNotFound:
		c.JSON(http.StatusNotFound, api.Error("order_not_found", "order not found"))
	default:
		c.JSON(http.StatusInternalServerError, api.Error("internal_error", "unexpected behaviour"))
	}

	c.JSON(http.StatusOK, gin.H{"message": "order confirmed"})
}

func (h *Hanlder) GetOrderByID(c *gin.Context) {
	var params = struct {
		ID string `uri:"id" binding:"required,uuid"`
	}{}

	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest,
			api.Error("invalid_request", err.Error()))
		return
	}

	orderID, err := uuid.Parse(params.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			api.Error("invalid_request", err.Error()))
		return
	}

	order, err := h.service.GetByID(c.Request.Context(), orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, api.Error("internal_error", err.Error()))
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
		Status:     model.OrderStatus(order.Status),
		TotalPrice: order.TotalPrice,
		Items:      items,
		CreatedAt:  order.CreatedAt,
		UpdatedAt:  order.UpdatedAt,
	}})
}

func (h *Hanlder) GetAllOrders(c *gin.Context) {
	var query api.PageQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, api.Error("invalid_request", err.Error()))
		return
	}

	if query.Limit == 0 {
		query.Limit = 10
	}

	if query.Page == 0 {
		query.Page = 1
	}

	orders, page, err := h.service.GetAll(c.Request.Context(), &query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, api.Error("internal_error", "unexpected behaviour"))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": orders,
		"meta": gin.H{
			"page":  page.Index,
			"limit": page.Limit,
			"total": page.Total,
		},
	})
}
