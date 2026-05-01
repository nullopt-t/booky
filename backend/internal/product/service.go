package product

import (
	"booky-backend/internal/domain"
	"booky-backend/internal/utils"
	"context"
)

type Service struct {
	repo Repository
}

func NewService(r Repository) *Service {
	return &Service{repo: r}
}

func (s *Service) Create(ctx context.Context, req CreateProductRequest) (*domain.Product, error) {
	return s.repo.Create(ctx, req)
}

func (s *Service) Update(ctx context.Context, id string, req UpdateProductRequest) (*domain.Product, error) {
	return s.repo.Update(ctx, id, req)
}

func (s *Service) GetAll(ctx context.Context, q utils.PaginationQuery) (*utils.PageResult[domain.Product], error) {
	return s.repo.GetAll(ctx, q)
}

func (s *Service) GetByID(ctx context.Context, id string) (*domain.Product, error) {
	return s.repo.GetByID(ctx, id)
}
