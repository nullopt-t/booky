package cart

import (
	"context"

	"booky-backend/internal/db"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CartRepository interface {
	Create(ctx context.Context, db db.DBQE, userID uuid.UUID) (*Cart, error)
	GetByUserID(ctx context.Context, db db.DBQE, userID uuid.UUID) (*Cart, error)
	Empty(ctx context.Context, db db.DBQE, userID uuid.UUID) error
	Save(ctx context.Context, db db.DBQE, cart *Cart) error
}

type CartService interface {
	GetCart(ctx context.Context, userID uuid.UUID) (*Cart, error)
	AddItem(ctx context.Context, userID uuid.UUID, req AddCartItemRequest) (*Cart, error)
	// RemoveItem(ctx context.Context, userID string, itemID string) (*Cart, error)
	EmptyCart(ctx context.Context, userID uuid.UUID) error
}

type CartHandler interface {
	GetCart(c *gin.Context)
	AddItem(c *gin.Context)
	EmptyCart(c *gin.Context)
}
