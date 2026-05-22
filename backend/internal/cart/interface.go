package cart

import (
	"context"

	"booky-backend/internal/model"
	"booky-backend/pkg/database"

	"github.com/gin-gonic/gin"
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

type CartService interface {
	GetCart(ctx context.Context, userID uuid.UUID) (*model.Cart, error)
	AddItem(ctx context.Context, userID uuid.UUID, req AddCartItemRequest) (*model.Cart, error)
	EmptyCart(ctx context.Context, userID uuid.UUID) error
}

type CartHandler interface {
	GetCart(c *gin.Context)
	AddItem(c *gin.Context)
	EmptyCart(c *gin.Context)
}
