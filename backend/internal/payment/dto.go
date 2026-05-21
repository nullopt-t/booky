package payment

import "time"

type PaymentResponse struct {
	ID          string     `json:"id"`
	OrderID     string     `json:"order_id"`
	Amount      float64    `json:"amount"`
	Status      string     `json:"status"`
	Provider    string     `json:"provider"`
	ProviderRef string     `json:"provider_ref"`
	PaidAt      *time.Time `json:"paid_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type CreatePaymentRequest struct {
	OrderID        string `json:"order_id" binding:"required"`
	IdempotencyKey string `json:"idempotency_key" binding:"required"`
	PaymentMethod  string `json:"payment_method,omitempty"`
}
