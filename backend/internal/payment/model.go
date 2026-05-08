package payment

import "time"

const (
	PaymentStatusPending    PaymentStatus = "pending"
	PaymentStatusProcessing PaymentStatus = "processing"
	PaymentStatusSucceeded  PaymentStatus = "succeeded"
	PaymentStatusFailed     PaymentStatus = "failed"
	PaymentStatusCancelled  PaymentStatus = "cancelled"
)

type Payment struct {
	ID          string
	OrderID     string
	Amount      int
	Status      PaymentStatus
	Provider    string
	ProviderRef string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
