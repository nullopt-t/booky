package cart

import (
	"booky-backend/internal/db"
	"booky-backend/internal/model"
	"context"

	"github.com/google/uuid"
)

type Service struct {
	tx          db.Runner
	cartRepo    CartRepository
	productRepo ProductRepository
}

func NewService(tx db.Runner, cartRepo CartRepository, productRepo ProductRepository) CartService {
	return &Service{tx, cartRepo, productRepo}
}

func (s *Service) getOrCreateCart(ctx context.Context, db db.DBQE, userID uuid.UUID) (*model.Cart, error) {
	cart, err := s.cartRepo.GetByUserID(ctx, db, userID)
	if err != nil {
		return s.cartRepo.Create(ctx, db, userID)
	}
	return cart, nil
}

func (s *Service) addOrUpdateItem(ctx context.Context, db db.DBQE, cart *model.Cart, req AddCartItemRequest) error {
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

	return s.cartRepo.Save(ctx, db, cart)
}

func (s *Service) GetCart(ctx context.Context, userID uuid.UUID) (*model.Cart, int, error) {
	var total int
	cart, err := s.getOrCreateCart(ctx, s.tx.DB(), userID)
	if err != nil {
		return nil, total, err
	}

	for _, item := range cart.Items {
		product, err := s.productRepo.GetByID(ctx, s.tx.DB(), item.ProductID)
		if err != nil {
			return nil, total, err
		}
		total += product.Price * item.Quantity
	}

	return cart, total, nil
}

func (s *Service) AddItem(ctx context.Context, userID uuid.UUID, req AddCartItemRequest) (*model.Cart, error) {

	var cart *model.Cart

	err := s.tx.WithTx(ctx, func(tx db.DBQE) error {

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
	return s.tx.WithTx(ctx, func(tx db.DBQE) error {
		return s.cartRepo.Empty(ctx, tx, userID)
	})
}
