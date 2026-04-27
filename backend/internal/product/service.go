package product

import (
	"booky-backend/internal/domain"
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

func (s *Service) List(ctx context.Context) ([]domain.Product, error) {
	return s.repo.List(ctx)
}

func (s *Service) GetByID(ctx context.Context, id string) (*domain.Product, error) {
	return s.repo.GetByID(ctx, id)
}
