package checkout

import (
	"booky-backend/internal/model"
	"booky-backend/pkg/database"
	"booky-backend/pkg/logger"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type Service struct {
	tx            database.Runner
	orderRepo     OrderRepository
	cartRepo      CartRepository
	productRepo   ProductRepository
	inventoryRepo InventoryRepository
}

func NewService(tx database.Runner, order OrderRepository, cart CartRepository, product ProductRepository, inventory InventoryRepository) CheckoutService {
	return &Service{
		tx:            tx,
		orderRepo:     order,
		cartRepo:      cart,
		productRepo:   product,
		inventoryRepo: inventory,
	}
}

func (s *Service) Checkout(ctx context.Context, userID uuid.UUID) error {
	return s.tx.WithTx(ctx, func(tx database.QueryExecutor) error {
		cart, err := s.cartRepo.GetByUserID(ctx, tx, userID)
		if err != nil {
			logger.Log(logger.ERROR, err.Error())
			return ErrNotFound
		}

		var totalPrice int
		var orderItems = make([]model.OrderItem, 0, len(cart.Items))
		for _, item := range cart.Items {
			p, err := s.productRepo.GetByID(ctx, tx, item.ProductID)
			if err != nil {
				logger.Log(logger.ERROR, err.Error())
				return ErrProductNotFound
			}

			availableQuantity, err := s.inventoryRepo.GetAvailable(ctx, tx, p.ID)
			if err != nil {
				logger.Log(logger.ERROR, err.Error())
				return err
			}

			if availableQuantity < item.Quantity {
				return ErrInsufficientQuantity
			}

			orderItems = append(orderItems, model.OrderItem{
				ProductID:     item.ProductID,
				Quantity:      item.Quantity,
				PurchasePrice: p.Price,
			})

			totalPrice += p.Price * item.Quantity
		}

		order, err := s.orderRepo.Create(ctx, tx, &model.Order{
			Items:      orderItems,
			TotalPrice: totalPrice,
		})
		if err != nil {
			return err
		}

		logger.Log(logger.INFO, fmt.Sprintf("%s order created", order.ID))

		return nil
	})
}
