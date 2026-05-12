package order

import (
	"booky-backend/internal/trans"
	"context"
)

type Service struct {
	repo OrderRepository
}

func NewService(r OrderRepository) OrderService {
	return &Service{repo: r}
}

func (s *Service) Create(ctx context.Context, req *CreateOrderRequest) (*Order, error) {
	if len(req.Items) <= 0 {
		return nil, ErrNoItems
	}

	order, err := s.repo.Create(ctx, req)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (s *Service) GetByID(ctx context.Context, orderID string) (*Order, error) {
	order, err := s.repo.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (s *Service) GetAll(ctx context.Context, q *trans.PaginationQuery) ([]*Order, *trans.Page, error) {
	orders, page, err := s.repo.GetAll(ctx, q)
	if err != nil {
		return nil, nil, err
	}
	return orders, page, nil
}

func (s *Service) Cancel(ctx context.Context, orderID string) error {
	return s.repo.TransitionStatus(ctx, orderID, OrderStatusPending, OrderStatusCancelled)
}

func (s *Service) Confirm(ctx context.Context, orderID string) error {
	return s.repo.TransitionStatus(ctx, orderID, OrderStatusPending, OrderStatusConfirmed)
}
