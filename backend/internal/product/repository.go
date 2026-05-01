package product

import (
	"booky-backend/internal/domain"
	"booky-backend/internal/utils"
	"context"
)

type Repository interface {
	Create(ctx context.Context, p CreateProductRequest) (*domain.Product, error)
	Update(ctx context.Context, id string, p UpdateProductRequest) (*domain.Product, error)
	GetByID(ctx context.Context, id string) (*domain.Product, error)
	GetAll(ctx context.Context, q utils.PaginationQuery) (*utils.PageResult[domain.Product], error)
}
