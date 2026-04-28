package product

import (
	"booky-backend/internal/db"
	"booky-backend/internal/domain"
	"context"
)

type Repository interface {
	Create(ctx context.Context, p CreateProductRequest) (*domain.Product, error)
	GetByID(ctx context.Context, id string) (*domain.Product, error)
	GetAll(ctx context.Context, q db.PaginationQuery) (*PaginatedProductsResponse, error)
}
