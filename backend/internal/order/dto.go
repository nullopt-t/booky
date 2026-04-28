package order

import (
	"booky-backend/internal/domain"
	"time"
)

type OrderItemResponse struct {
	ProductID     string `json:"product_id"`
	Quantity      int    `json:"quantity"`
	PurchasePrice int    `json:"purchase_price"`
}

type OrderResponse struct {
	ID         string              `json:"id"`
	Status     string              `json:"status"`
	TotalPrice int64               `json:"total_price"`
	Items      []OrderItemResponse `json:"items"`
	CreatedAt  time.Time           `json:"created_at"`
	UpdatedAt  time.Time           `json:"updated_at"`
}

type CreateOrderItemRequest struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type CreateOrderRequest struct {
	Items []CreateOrderItemRequest `json:"items"`
}

type CreateOrderResponse struct {
	ID         string              `json:"id"`
	Status     domain.OrderStatus  `json:"status"`
	Items      []OrderItemResponse `json:"items"`
	TotalPrice int                 `json:"total_price"`
	ItemsCount int                 `json:"items_count"`
	CreatedAt  time.Time           `json:"created_at"`
}
