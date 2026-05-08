package product

import (
	"booky-backend/internal/trans"
	"context"
)

type ProductRepository interface {
	Create(ctx context.Context, p CreateProductRequest) (*Product, error)
	Update(ctx context.Context, id string, p UpdateProductRequest) (*Product, error)
	GetByID(ctx context.Context, id string) (*Product, error)
	GetAll(ctx context.Context, q trans.PaginationQuery) ([]Product, *trans.Page, error)
}
