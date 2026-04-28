package order

import (
	"booky-backend/internal/db"
	"context"
)

type Service struct {
	repo Repository
}

func NewService(r Repository) *Service {
	return &Service{repo: r}
}

func (s *Service) Create(ctx context.Context, order CreateOrderRequest) (*CreateOrderResponse, error) {
	return s.repo.Create(ctx, order)
}

func (s *Service) GetAll(ctx context.Context, q db.PaginationQuery) ([]*OrderResponse, error) {
	return s.repo.GetAll(ctx, q)
}
