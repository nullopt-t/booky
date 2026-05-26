package order

import (
	"booky-backend/internal/model"
	"booky-backend/pkg/api"
	"booky-backend/pkg/database"
	"context"

	"github.com/google/uuid"
)

type OrderRepository interface {
	Create(ctx context.Context, db database.QueryExecutor, order model.Order) (*model.Order, error)

	GetByID(ctx context.Context, db database.QueryExecutor, orderID uuid.UUID) (*model.Order, error)

	GetAll(ctx context.Context, db database.QueryExecutor, q *api.PageQuery) ([]*model.Order, *api.Page, error)

	TransitionStatus(ctx context.Context, db database.QueryExecutor, orderID uuid.UUID, from, to model.OrderStatus) error

	UpdateTotalPrice(ctx context.Context, db database.QueryExecutor, orderID uuid.UUID, total int) error
}

type Service struct {
	dbExecuter database.Runner
	repo       OrderRepository
}

func NewService(dbExecuter database.Runner, r OrderRepository) OrderService {
	return &Service{dbExecuter: dbExecuter, repo: r}
}

func (s *Service) GetByID(ctx context.Context, orderID uuid.UUID) (*model.Order, error) {
	var order *model.Order
	err := s.dbExecuter.WithDB(ctx, func(db database.QueryExecutor) error {
		var err error
		order, err = s.repo.GetByID(ctx, db, orderID)
		return err
	})
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (s *Service) GetAll(ctx context.Context, q *api.PageQuery) ([]*model.Order, *api.Page, error) {
	var orders []*model.Order
	var page *api.Page
	err := s.dbExecuter.WithDB(ctx, func(db database.QueryExecutor) error {
		var err error
		orders, page, err = s.repo.GetAll(ctx, db, q)
		return err
	})
	if err != nil {
		return nil, nil, err
	}
	return orders, page, nil
}

func (s *Service) Cancel(ctx context.Context, orderID uuid.UUID) error {
	return s.dbExecuter.WithTx(ctx, func(tx database.QueryExecutor) error {
		return s.repo.TransitionStatus(ctx, tx, orderID, model.OrderStatusPending, model.OrderStatusCancelled)
	})
}

func (s *Service) Confirm(ctx context.Context, orderID uuid.UUID) error {
	return s.dbExecuter.WithTx(ctx, func(tx database.QueryExecutor) error {
		return s.repo.TransitionStatus(ctx, tx, orderID, model.OrderStatusPending, model.OrderStatusConfirmed)
	})
}
