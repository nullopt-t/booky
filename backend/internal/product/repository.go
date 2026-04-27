package product

import (
	"booky-backend/internal/domain"
	"context"
)

type Repository interface {
	Create(ctx context.Context, p CreateProductRequest) (*domain.Product, error)
	GetByID(ctx context.Context, id string) (*domain.Product, error)
	List(ctx context.Context) ([]domain.Product, error)
}
