package cart

import (
	"context"

	"booky-backend/pkg/database"
	"booky-backend/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProductRepository interface {
	GetByID(ctx context.Context, db database.DBQE, productID uuid.UUID) (*model.Product, error)
}

type CartRepository interface {
	Create(ctx context.Context, db database.DBQE, userID uuid.UUID) (*model.Cart, error)
	GetByUserID(ctx context.Context, db database.DBQE, userID uuid.UUID) (*model.Cart, error)
	Empty(ctx context.Context, db database.DBQE, userID uuid.UUID) error
	Save(ctx context.Context, db database.DBQE, cart *model.Cart) error
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
