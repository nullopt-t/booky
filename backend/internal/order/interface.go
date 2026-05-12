package order

import (
	"booky-backend/internal/trans"
	"context"

	"github.com/gin-gonic/gin"
)

type OrderRepository interface {
	// Create creates a new order with items already included in the request context
	Create(ctx context.Context, order *CreateOrderRequest) (*Order, error)

	// GetByID returns a single order with its items
	GetByID(ctx context.Context, orderID string) (*Order, error)

	// GetAll returns paginated orders list
	GetAll(ctx context.Context, q *trans.PaginationQuery) ([]*Order, *trans.Page, error)

	TransitionStatus(ctx context.Context, orderID string, from, to OrderStatus) error

	// LockForUpdate locks order row for safe transactional operations
	// (used inside payment/webhook flows)
	LockForUpdate(ctx context.Context, orderID string) (*Order, error)

	// UpdateTotalPrice recalculates or sets total price
	UpdateTotalPrice(ctx context.Context, orderID string, total int) error
}

type OrderService interface {
	Create(ctx context.Context, req *CreateOrderRequest) (*Order, error)
	GetByID(ctx context.Context, orderID string) (*Order, error)
	GetAll(ctx context.Context, q *trans.PaginationQuery) ([]*Order, *trans.Page, error)
	Cancel(ctx context.Context, orderID string) error
	Confirm(ctx context.Context, orderID string) error
}

type OrderHandler interface {
	CreateOrder(c *gin.Context)
	GetOrderByID(c *gin.Context)
	GetAllOrders(c *gin.Context)
	CancelOrder(c *gin.Context)
	ConfirmOrder(c *gin.Context)
}
