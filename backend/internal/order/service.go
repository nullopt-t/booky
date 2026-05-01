package order

import (
	"booky-backend/internal/domain"
	"booky-backend/internal/utils"
	"context"
	"log"
)

type Service struct {
	repo Repository
}

func NewService(r Repository) *Service {
	return &Service{repo: r}
}

func (s *Service) Create(ctx context.Context, order CreateOrderRequest) (*domain.Order, error) {
	createdOrder, err := s.repo.Create(ctx, order)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return createdOrder, nil
}

func (s *Service) Cancel(ctx context.Context, orderID string) error {
	err := s.repo.Cancel(ctx, orderID)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (s *Service) Confirm(ctx context.Context, orderID string) error {
	err := s.repo.Confirm(ctx, orderID)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (s *Service) GetByID(ctx context.Context, id string) (*domain.Order, error) {
	order, err := s.repo.GetByID(ctx, id)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return order, nil
}

func (s *Service) GetAll(ctx context.Context, q utils.PaginationQuery) ([]*OrderResponse, error) {
	return s.repo.GetAll(ctx, q)
}
