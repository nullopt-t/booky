package cart

import "time"

type AddCartItemRequest struct {
	ItemID   string `json:"item_id"`
	Quantity int    `json:"quantity"`
}

type CartItemResponse struct {
	ItemID   string `json:"item_id"`
	Quantity int    `json:"quantity"`
}

type CartResponse struct {
	ID        string             `json:"id"`
	Items     []CartItemResponse `json:"items"`
	UpdatedAt time.Time          `json:"updated_at"`
}
