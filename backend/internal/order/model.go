package order

import "time"

type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusPaid       OrderStatus = "paid"
	OrderStatusConfirmed  OrderStatus = "confirmed"
	OrderStatusProcessing OrderStatus = "processing"
	OrderStatusShipped    OrderStatus = "shipped"
	OrderStatusDelivered  OrderStatus = "delivered"
	OrderStatusCancelled  OrderStatus = "cancelled"
	OrderStatusRefuneded  OrderStatus = "refuneded"
)

type OrderItem struct {
	ProductID        string
	Quantity      int
	PurchasePrice int
}

type Order struct {
	// order data
	ID         string
	Status     OrderStatus
	TotalPrice int

	// items data
	Items []OrderItem

	// timestamps
	CreatedAt time.Time
	UpdatedAt time.Time
}
