package order

import (
	"time"
)



type OrderItemResponse struct {
	ProductID     string `json:"product_id"`
	Quantity      int    `json:"quantity"`
	PurchasePrice int    `json:"purchase_price"`
}

type OrderResponse struct {
	ID         string              `json:"id"`
	Status     OrderStatus         `json:"status"`
	TotalPrice int                 `json:"total_price"`
	Items      []OrderItemResponse `json:"items"`
	CreatedAt  time.Time           `json:"created_at"`
	UpdatedAt  time.Time           `json:"updated_at"`
}

type CreateOrderItemRequest struct {
	ProductID     string `json:"product_id" binding:"required,uuid"`
	Quantity      int    `json:"quantity" binding:"required,min=1,max=100"`
}

type CreateOrderRequest struct {
	Items []CreateOrderItemRequest `json:"items" binding:"required,min=1,dive"`
}

type CreateOrderResponse struct {
	ID         string              `json:"id"`
	Status     string              `json:"status"`
	Items      []OrderItemResponse `json:"items"`
	TotalPrice int                 `json:"total_price"`
	CreatedAt  time.Time           `json:"created_at"`
}
