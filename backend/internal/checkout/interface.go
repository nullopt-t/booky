package checkout

import (
	"booky-backend/internal/model"
	"booky-backend/pkg/database"
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CartRepository interface {
	GetByUserID(ctx context.Context, db database.QueryExecutor, userID uuid.UUID) (*model.Cart, error)
	Empty(ctx context.Context, db database.QueryExecutor, userID uuid.UUID) error
	Save(ctx context.Context, db database.QueryExecutor, cart *model.Cart) error
}

type ProductRepository interface {
	GetByID(ctx context.Context, db database.QueryExecutor, productID uuid.UUID) (*model.Product, error)
}

type OrderRepository interface {
	Create(ctx context.Context, db database.QueryExecutor, order *model.Order) (*model.Order, error)
}

type InventoryRepository interface {
	Reserve(ctx context.Context, db database.QueryExecutor, productID uuid.UUID, quantity int) error
	Release(ctx context.Context, pdb database.QueryExecutor, roductID uuid.UUID, quantity int) error
	GetAvailable(ctx context.Context, db database.QueryExecutor, productID uuid.UUID) (int, error)
	GetReserved(ctx context.Context, db database.QueryExecutor, productID uuid.UUID) (int, error)
}

type CheckoutService interface {
	Checkout(ctx context.Context, userID uuid.UUID) error
}

type CheckoutHandler interface {
	HandleCheckout(c *gin.Context)
}
