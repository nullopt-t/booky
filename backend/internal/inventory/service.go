package inventory

import (
	"booky-backend/pkg/database"
	"context"

	"github.com/google/uuid"
)

type Service struct {
	tx   database.Runner
	repo InventoryRepository
}

func NewService(tx database.Runner, repo InventoryRepository) *Service {
	return &Service{
		tx:   tx,
		repo: repo,
	}
}

func (s *Service) GetAvailable(ctx context.Context, productID uuid.UUID) (int, error) {
	count, err := s.repo.GetAvailable(ctx, s.tx.DB(), productID)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (s *Service) GetReserved(ctx context.Context, productID uuid.UUID) (int, error) {
	count, err := s.repo.GetReserved(ctx, s.tx.DB(), productID)
	if err != nil {
		return 0, err
	}
	return count, nil
}
