package cart

import (
	"time"

	"github.com/google/uuid"
)

type AddCartItemRequest struct {
	ItemID   uuid.UUID `json:"item_id"`
	Quantity int       `json:"quantity"`
}

type CartItemResponse struct {
	ItemID   uuid.UUID `json:"item_id"`
	Quantity int       `json:"quantity"`
}

type CartResponse struct {
	ID        uuid.UUID          `json:"id"`
	Total     *int               `json:"total,omitempty"`
	Items     []CartItemResponse `json:"items"`
	UpdatedAt time.Time          `json:"updated_at"`
}
