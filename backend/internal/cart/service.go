package cart

import (
	"booky-backend/internal/db"
	"context"

	"github.com/google/uuid"
)

type Service struct {
	tx   *db.TxRunner
	repo CartRepository
}

func NewService(repo CartRepository, tx *db.TxRunner) CartService {
	return &Service{
		repo: repo,
		tx:   tx,
	}
}

func (s *Service) getOrCreateCart(ctx context.Context, db db.DBQE, userID uuid.UUID) (*Cart, error) {
	cart, err := s.repo.GetByUserID(ctx, db, userID)
	if err != nil {
		if err == ErrCartNotFound {
			cart, err = s.repo.Create(ctx, db, userID)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return cart, nil
}

func (s *Service) addOrUpdateItem(ctx context.Context, db db.DBQE, cart *Cart, req AddCartItemRequest) error {
	found := false
	for i, item := range cart.Items {
		if item.ProductID == req.ItemID {
			cart.Items[i].Quantity += req.Quantity
			found = true
			break
		}
	}
	if !found {
		cart.Items = append(cart.Items, CartItem{
			ProductID: req.ItemID,
			Quantity:  req.Quantity,
		})
	}

	return s.repo.Save(ctx, db, cart)
}

func (s *Service) GetCart(ctx context.Context, userID uuid.UUID) (*Cart, error) {
	return s.getOrCreateCart(ctx, s.tx.DB(), userID)
}

func (s *Service) AddItem(ctx context.Context, userID uuid.UUID, req AddCartItemRequest) (*Cart, error) {

	var cart *Cart

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
		return s.repo.Empty(ctx, tx, userID)
	})
}
