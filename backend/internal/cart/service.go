package cart

import (
	"booky-backend/internal/model"
	"booky-backend/pkg/database"
	"context"
	"errors"

	"github.com/google/uuid"
)

type Service struct {
	tx          database.Runner
	cartRepo    CartRepository
	productRepo ProductRepository
}

func NewService(tx database.Runner, cartRepo CartRepository, productRepo ProductRepository) CartService {
	return &Service{tx, cartRepo, productRepo}
}

func (s *Service) getOrCreateCart(ctx context.Context, qe database.QueryExecutor, userID uuid.UUID) (*model.Cart, error) {
	cart, err := s.cartRepo.GetByUserID(ctx, qe, userID)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return s.cartRepo.Create(ctx, qe, userID)
		}
		return nil, err
	}
	return cart, nil
}

func (s *Service) addOrUpdateItem(ctx context.Context, qe database.QueryExecutor, cart *model.Cart, req AddCartItemRequest) error {
	found := false
	for i, item := range cart.Items {
		if item.ProductID == req.ProductID {
			cart.Items[i].Quantity += req.Quantity
			found = true
			break
		}
	}
	if !found {
		cart.Items = append(cart.Items, model.CartItem{
			ProductID: req.ProductID,
			Quantity:  req.Quantity,
		})
	}

	if err := s.cartRepo.Save(ctx, qe, cart); err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return ErrCartNotFound
		}
		return err
	}
	return nil
}

func (s *Service) GetCart(ctx context.Context, userID uuid.UUID) (*model.Cart, int, error) {
	var total int
	cart, err := s.getOrCreateCart(ctx, s.tx.DB(), userID)
	if err != nil {
		if errors.Is(err, database.ErrConflict) {
			return nil, total, ErrCartAlreadyExist
		}
		return nil, total, err
	}

	for _, item := range cart.Items {
		p, err := s.productRepo.GetByID(ctx, s.tx.DB(), item.ProductID)
		if err != nil {
			if errors.Is(err, database.ErrNotFound) {
				return nil, total, ErrProductNotFound
			}
			return nil, total, err
		}
		total += p.Price * item.Quantity
	}

	return cart, total, nil
}

func (s *Service) AddItem(ctx context.Context, userID uuid.UUID, req AddCartItemRequest) (*model.Cart, error) {

	var cart *model.Cart

	err := s.tx.WithTx(ctx, func(tx database.QueryExecutor) error {

		c, err := s.getOrCreateCart(ctx, tx, userID)
		if err != nil {
			return err
		}

		if err := s.addOrUpdateItem(ctx, tx, c, req); err != nil {
			return err
		}

		cart = c
		return nil
	})

	return cart, err
}

func (s *Service) EmptyCart(ctx context.Context, userID uuid.UUID) error {
	return s.tx.WithTx(ctx, func(tx database.QueryExecutor) error {
		return s.cartRepo.Empty(ctx, tx, userID)
	})
}
