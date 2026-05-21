package order

import (
	"booky-backend/internal/db"
	"booky-backend/internal/model"
	"booky-backend/internal/trans"
	"context"

	"github.com/google/uuid"
)

type Service struct {
	tx   db.Runner
	repo OrderRepository
}

func NewService(tx db.Runner, r OrderRepository) OrderService {
	return &Service{tx: tx, repo: r}
}

func (s *Service) GetByID(ctx context.Context, orderID uuid.UUID) (*model.Order, error) {
	order, err := s.repo.GetByID(ctx, s.tx.DB(), orderID)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (s *Service) GetAll(ctx context.Context, q *trans.PaginationQuery) ([]*model.Order, *trans.Page, error) {
	orders, page, err := s.repo.GetAll(ctx, s.tx.DB(), q)
	if err != nil {
		return nil, nil, err
	}
	return orders, page, nil
}

func (s *Service) Cancel(ctx context.Context, orderID uuid.UUID) error {
	return s.tx.WithTx(ctx, func(tx db.DBQE) error {
		return s.repo.TransitionStatus(ctx, s.tx.DB(), orderID, model.OrderStatusPending, model.OrderStatusCancelled)
	})
}

func (s *Service) Confirm(ctx context.Context, orderID uuid.UUID) error {
	return s.tx.WithTx(ctx, func(tx db.DBQE) error {
		return s.repo.TransitionStatus(ctx, s.tx.DB(), orderID, model.OrderStatusPending, model.OrderStatusConfirmed)
	})
}
