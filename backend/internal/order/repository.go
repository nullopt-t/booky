package order

import (
	"booky-backend/internal/trans"
	"context"
)

type OrderRepository interface {
	Create(ctx context.Context, order CreateOrderRequest) (*Order, error)
	Cancel(ctx context.Context, orderID string) error
	Confirm(ctx context.Context, orderID string) error
	GetByID(ctx context.Context, id string) (*Order, error)
	GetAll(ctx context.Context, q trans.PaginationQuery) ([]Order, *trans.Page, error)
}
