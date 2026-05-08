package product

import (
	"booky-backend/internal/trans"
	"context"
)

type Service struct {
	repo ProductRepository
}

func NewService(r ProductRepository) *Service {
	return &Service{repo: r}
}

func (s *Service) Create(ctx context.Context, req CreateProductRequest) (*Product, error) {
	return s.repo.Create(ctx, req)
}

func (s *Service) Update(ctx context.Context, id string, req UpdateProductRequest) (*Product, error) {
	return s.repo.Update(ctx, id, req)
}

func (s *Service) GetAll(ctx context.Context, q trans.PaginationQuery) ([]Product, *trans.Page, error) {
	return s.repo.GetAll(ctx, q)
}

func (s *Service) GetByID(ctx context.Context, id string) (*Product, error) {
	return s.repo.GetByID(ctx, id)
}
