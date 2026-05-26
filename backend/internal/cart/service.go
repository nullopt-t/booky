package cart

import (
	"booky-backend/internal/model"
	"booky-backend/pkg/database"
	"context"
	"errors"

	"github.com/google/uuid"
)

type ProductRepository interface {
	GetByID(ctx context.Context, db database.QueryExecutor, productID uuid.UUID) (*model.Product, error)
}
type CartRepository interface {
	Create(ctx context.Context, db database.QueryExecutor, userID uuid.UUID) (*model.Cart, error)
	GetByUserID(ctx context.Context, db database.QueryExecutor, userID uuid.UUID) (*model.Cart, error)
	Empty(ctx context.Context, db database.QueryExecutor, userID uuid.UUID) error
	Save(ctx context.Context, db database.QueryExecutor, cart *model.Cart) error
}
type Service struct {
	dbExecutor  database.Runner
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
		if item.ProductID.String() == req.ItemID.String() {
			cart.Items[i].Quantity += req.Quantity
			found = true
			break
		}
	}
	if !found {
		cart.Items = append(cart.Items, model.CartItem{
			ProductID: req.ItemID,
			Quantity:  req.Quantity,
		})
	}

	if err := s.cartRepo.Save(ctx, qe, cart); err != nil {
		return err
	}
	return nil
}

func (s *Service) GetCart(ctx context.Context, userID uuid.UUID) (*model.Cart, int, error) {
	var total int
	var cart *model.Cart
	err := s.dbExecutor.WithDB(ctx, func(db database.QueryExecutor) error {
		var err error
		cart, err = s.getOrCreateCart(ctx, db, userID)
		if err != nil {
			return err
		}
		for _, item := range cart.Items {
			p, err := s.productRepo.GetByID(ctx, db, item.ProductID)
			if err != nil {
				return err
			}
			total += p.Price * item.Quantity
		}
		return nil
	})
	if err != nil {
		return nil, 0, err
	}
	return cart, total, nil
}

func (s *Service) AddItem(ctx context.Context, userID uuid.UUID, req AddCartItemRequest) (*model.Cart, error) {

	var cart *model.Cart

	err := s.dbExecutor.WithTx(ctx, func(tx database.QueryExecutor) error {

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
	return s.dbExecutor.WithTx(ctx, func(tx database.QueryExecutor) error {
		return s.cartRepo.Empty(ctx, tx, userID)
	})
}
