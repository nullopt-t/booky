package payment

import (
	"context"
)

type PaymentRepository interface {
	// GetByID returns a payment by its ID.
	GetByID(ctx context.Context, id string) (*Payment, error)

	// GetByProviderRef returns a payment by its provider ref.
	GetByProviderRef(ctx context.Context, provider string, providerRef string) (*Payment, error)

	// TransitionStatus transitions a payment's status from the specified old status to the specified new status.
	TransitionStatus(ctx context.Context, id string, oldStatus PaymentStatus, newStatus PaymentStatus) error

	TransitionStatusByProviderRef(ctx context.Context, provider string, providerRef string, newStatus PaymentStatus) error
	// Create saves a new payment.
	Create(ctx context.Context, orderID string, provider string, providerRef string) (*Payment, error)
}
