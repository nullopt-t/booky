package domain

import "time"

type OrderItem struct {
	ProductID     string
	Quantity      int
	PurchasePrice int
}

type Order struct {
	// order data
	ID         string
	Status     string
	TotalPrice int

	// items data
	Items []OrderItem

	// timestamps
	CreatedAt time.Time
	UpdatedAt time.Time
}
