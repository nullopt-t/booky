package domain

import "time"

type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusConfirmed  OrderStatus = "confirmed"
	OrderStatusPaid       OrderStatus = "paid"
	OrderStatusProcessing OrderStatus = "processing"
	OrderStatusShipped    OrderStatus = "shipped"
	OrderStatusDelivered  OrderStatus = "delivered"
	OrderStatusCancelled  OrderStatus = "cancelled"
	OrderStatusRefuneded  OrderStatus = "refuneded"
)

type OrderItem struct {
	ProductID     string `json:"product_id" binding:"required"`
	Quantity      int    `json:"quantity" binding:"required,min=1"`
	PurchasePrice int    `json:"purchase_price,omitempty"`
}

type Order struct {
	// order data
	ID         string      `json:"id"`
	Status     OrderStatus `json:"status"`
	TotalPrice int         `json:"total_price"`

	// items data
	Items []OrderItem `json:"items"`

	// timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
