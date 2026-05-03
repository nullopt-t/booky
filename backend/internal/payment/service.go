package payment

import (
	"booky-backend/internal/domain"
	"context"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return Service{
		repo: repo,
	}
}

func (s *Service) CreatePayment(ctx context.Context, orderID string) error {
	return s.repo.Create(ctx, orderID)
}

func (s *Service) GetPaymentByID(ctx context.Context, orderID string) (*domain.Payment, error) {
	return s.repo.FindByID(orderID)
}
