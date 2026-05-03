package payment

import "time"

type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusPaid      PaymentStatus = "paid"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusCancelled PaymentStatus = "cancelled"
)

type Payment struct {
	ID          string
	OrderID     string
	Amount      int
	Status      PaymentStatus
	Provider    string
	ProviderRef string
	PaidAt      *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

