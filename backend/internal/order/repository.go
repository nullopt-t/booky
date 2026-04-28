package order

import (
	"booky-backend/internal/db"
	"context"
)

type Repository interface {
	Create(ctx context.Context, order CreateOrderRequest) (*CreateOrderResponse, error)
	// Update(order *domain.Order) error
	// Delete(id string) error
	// GetByID(id string) (*domain.Order, error)
	GetAll(ctx context.Context, q db.PaginationQuery) ([]*OrderResponse, error)
}
