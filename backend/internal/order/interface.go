package order

import (
	"booky-backend/pkg/database"
	"booky-backend/internal/model"
	"booky-backend/internal/trans"
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type OrderRepository interface {
	Create(ctx context.Context, db database.DBQE, order model.Order) (*model.Order, error)

	GetByID(ctx context.Context, db database.DBQE, orderID uuid.UUID) (*model.Order, error)

	GetAll(ctx context.Context, db database.DBQE, q *trans.PaginationQuery) ([]*model.Order, *trans.Page, error)

	TransitionStatus(ctx context.Context, db database.DBQE, orderID uuid.UUID, from, to model.OrderStatus) error

	UpdateTotalPrice(ctx context.Context, db database.DBQE, orderID uuid.UUID, total int) error
}

type OrderService interface {
	GetByID(ctx context.Context, orderID uuid.UUID) (*model.Order, error)
	GetAll(ctx context.Context, q *trans.PaginationQuery) ([]*model.Order, *trans.Page, error)
	Cancel(ctx context.Context, orderID uuid.UUID) error
	Confirm(ctx context.Context, orderID uuid.UUID) error
}

type OrderHandler interface {
	GetOrderByID(c *gin.Context)
	GetAllOrders(c *gin.Context)
	CancelOrder(c *gin.Context)
	ConfirmOrder(c *gin.Context)
}
