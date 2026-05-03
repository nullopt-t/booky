package domain

import "time"

type Payment struct {
	ID          string
	OrderID     string
	Amount      int
	Status      string
	Provider    string
	ProviderRef string
	PaidAt      *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
