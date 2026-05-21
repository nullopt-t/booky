package payment

import (
	"context"

	"github.com/google/uuid"
)

type PaymentRepository interface {
	// Create inserts a new payment (initial state: pending)
	Create(ctx context.Context, req CreatePaymentRequest) (*Payment, error)

	// GetByID fetches payment by internal ID
	GetByID(ctx context.Context, paymentID uuid.UUID) (*Payment, error)
	
	// GetByProviderRef fetches payment using external provider reference
	GetByProviderRef(ctx context.Context, provider string, providerRef string) (*Payment, error)

	// TransitionStatus atomically updates payment status (idempotent)
	TransitionStatus(ctx context.Context, paymentID uuid.UUID, oldStatus, newStatus PaymentStatus) error

	// TransitionStatusByProviderRef is used by webhook (idempotent external update)
	TransitionStatusByProviderRef(ctx context.Context, provider string, providerRef string, oldStatus, newStatus PaymentStatus) error

	// GetByOrderID fetches all payments for a given order
	GetByOrderID(ctx context.Context, orderID uuid.UUID) ([]Payment, error)
}

type PaymentService interface {
	// CreatePayment creates a new payment and returns the payment URL
	CreatePayment(ctx context.Context, req CreatePaymentRequest) (*PaymentResponse, error)

	// GetPayment retrieves a payment by ID
	GetPayment(ctx context.Context, paymentID uuid.UUID) (*Payment, error)

	// GetPaymentByProviderRef retrieves a payment by provider reference
	GetPaymentByProviderRef(ctx context.Context, provider string, providerRef string) (*Payment, error)

	// TransitionStatus transitions a payment to a new status (idempotent)
	TransitionStatus(ctx context.Context, paymentID uuid.UUID, oldStatus, newStatus PaymentStatus) error

	// TransitionStatusByProviderRef transitions a payment by provider reference (idempotent)
	TransitionStatusByProviderRef(ctx context.Context, provider string, providerRef string, oldStatus, newStatus PaymentStatus) error

	// GetByOrderID retrieves all payments for a given order
	GetByOrderID(ctx context.Context, orderID uuid.UUID) ([]Payment, error)
}
