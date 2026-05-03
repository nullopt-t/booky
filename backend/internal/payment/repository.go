package payment

import (
	"booky-backend/internal/domain"
	"context"
)

type Repository interface {
	// FindByID returns a payment by its ID.
	FindByID(id string) (*domain.Payment, error)

	// Create saves a new payment.
	Create(ctx context.Context, orderID string) error
}
