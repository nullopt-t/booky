package order

import (
	"booky-backend/internal/db"
	"booky-backend/internal/trans"
	"context"

	"github.com/google/uuid"
)

type Service struct {
	tx   *db.TxRunner
	repo OrderRepository
}

func NewService(tx *db.TxRunner, r OrderRepository) OrderService {
	return &Service{tx: tx, repo: r}
}

func (s *Service) Create(ctx context.Context, req *CreateOrderRequest) (*Order, error) {
	if len(req.Items) <= 0 {
		return nil, ErrNoItems
	}

	var order *Order
	err := s.tx.WithTx(ctx, func(tx db.DBQE) error {
		var err error
		order, err = s.repo.Create(ctx, tx, req)
		return err
	})

	return order, err
}

func (s *Service) GetByID(ctx context.Context, orderID uuid.UUID) (*Order, error) {
	order, err := s.repo.GetByID(ctx, s.tx.DB(), orderID)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (s *Service) GetAll(ctx context.Context, q *trans.PaginationQuery) ([]*Order, *trans.Page, error) {
	orders, page, err := s.repo.GetAll(ctx, s.tx.DB(), q)
	if err != nil {
		return nil, nil, err
	}
	return orders, page, nil
}

func (s *Service) Cancel(ctx context.Context, orderID uuid.UUID) error {
	return s.tx.WithTx(ctx, func(tx db.DBQE) error {
		return s.repo.TransitionStatus(ctx, s.tx.DB(), orderID, OrderStatusPending, OrderStatusCancelled)
	})
}

func (s *Service) Confirm(ctx context.Context, orderID uuid.UUID) error {
	return s.tx.WithTx(ctx, func(tx db.DBQE) error {
		return s.repo.TransitionStatus(ctx, s.tx.DB(), orderID, OrderStatusPending, OrderStatusConfirmed)
	})
}
