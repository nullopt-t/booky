package order

import (
	"booky-backend/internal/trans"
	"context"
)

type Service struct {
	repo OrderRepository
}

func NewService(r OrderRepository) *Service {
	return &Service{repo: r}
}

func (s *Service) Create(ctx context.Context, order CreateOrderRequest) (*Order, error) {
	if len(order.Items) == 0 {
		return nil, ErrNoItems
	}

	createdOrder, err := s.repo.Create(ctx, order)
	if err != nil {
		return nil, err
	}
	return createdOrder, nil
}

func (s *Service) Cancel(ctx context.Context, orderID string) error {
	err := s.repo.Cancel(ctx, orderID)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) Confirm(ctx context.Context, orderID string) error {
	err := s.repo.Confirm(ctx, orderID)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) GetByID(ctx context.Context, id string) (*Order, error) {
	order, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (s *Service) GetAll(ctx context.Context, q trans.PaginationQuery) ([]Order, *trans.Page, error) {
	return s.repo.GetAll(ctx, q)
}
