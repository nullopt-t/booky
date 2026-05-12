package cart

import (
	"context"
)

type Service struct {
	repo CartRepository
}

func NewService(repo CartRepository) CartService {
	return &Service{
		repo: repo,
	}
}

func (s *Service) getOrCreateCart(ctx context.Context, userID string) (*Cart, error) {
	cart, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		if err == ErrCartNotFound {
			cart, err = s.repo.Create(ctx, userID)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return cart, nil
}

func (s *Service) addOrUpdateItem(ctx context.Context, cart *Cart, req AddCartItemRequest) error {
	found := false
	for i, item := range cart.Items {
		if item.ItemID == req.ItemID {
			cart.Items[i].Quantity += req.Quantity
			found = true
			break
		}
	}
	if !found {
		cart.Items = append(cart.Items, CartItem{
			ItemID:   req.ItemID,
			Quantity: req.Quantity,
		})
	}

	return s.repo.Save(ctx, cart)
}

func (s *Service) GetCart(ctx context.Context, userID string) (*Cart, error) {
	return s.getOrCreateCart(ctx, userID)
}

func (s *Service) AddItem(ctx context.Context, userID string, req AddCartItemRequest) (*Cart, error) {
	cart, err := s.getOrCreateCart(ctx, userID)
	if err != nil {
		return nil, err
	}

	if err := s.addOrUpdateItem(ctx, cart, req); err != nil {
		return nil, err
	}

	return cart, nil
}

func (s *Service) Empty(ctx context.Context, userID string) error {
	return s.repo.Empty(ctx, userID)
}
