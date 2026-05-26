package inventory

import (
	"booky-backend/pkg/database"
	"context"

	"github.com/google/uuid"
)

type InventoryRepository interface {
	Reserve(ctx context.Context, qe database.QueryExecutor, productID uuid.UUID, quantity int) error
	Release(ctx context.Context, qe database.QueryExecutor, roductID uuid.UUID, quantity int) error
	GetAvailable(ctx context.Context, qe database.QueryExecutor, productID uuid.UUID) (int, error)
	GetReserved(ctx context.Context, qe database.QueryExecutor, productID uuid.UUID) (int, error)
}

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
	var count int
	err := s.tx.WithDB(ctx, func(db database.QueryExecutor) error {
		var err error
		count, err = s.repo.GetAvailable(ctx, db, productID)
		return err
	})
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (s *Service) GetReserved(ctx context.Context, productID uuid.UUID) (int, error) {
	var count int
	err := s.tx.WithDB(ctx, func(db database.QueryExecutor) error {
		var err error
		count, err = s.repo.GetReserved(ctx, db, productID)
		return err
	})
	if err != nil {
		return 0, err
	}
	return count, nil
}
