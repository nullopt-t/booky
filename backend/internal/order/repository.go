package order

import (
	"booky-backend/internal/domain"
	"booky-backend/internal/utils"
	"context"
)

type Repository interface {
	Create(ctx context.Context, order CreateOrderRequest) (*domain.Order, error)
	Cancel(ctx context.Context, orderID string) error
	Confirm(ctx context.Context, orderID string) error
	GetByID(ctx context.Context, id string) (*domain.Order, error)
	GetAll(ctx context.Context, q utils.PaginationQuery) ([]*OrderResponse, error)
}
