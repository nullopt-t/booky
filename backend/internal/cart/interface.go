package cart

import (
	"context"

	"github.com/gin-gonic/gin"
)

type CartRepository interface {
	Create(ctx context.Context, userID string) (*Cart, error)
	GetByUserID(ctx context.Context, userID string) (*Cart, error)
	Empty(ctx context.Context, userID string) error
	Save(ctx context.Context, cart *Cart) error
}

type CartService interface {
	GetCart(ctx context.Context, userID string) (*Cart, error)
	AddItem(ctx context.Context, userID string, req AddCartItemRequest) (*Cart, error)
	// RemoveItem(ctx context.Context, userID string, itemID string) (*Cart, error)
	Empty(ctx context.Context, userID string) error
}

type CartHandler interface {
	GetCart(c *gin.Context)
	AddItem(c *gin.Context)
	Empty(c *gin.Context)
}
