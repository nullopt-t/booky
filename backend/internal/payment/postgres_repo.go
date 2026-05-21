package payment

import (
	"booky-backend/internal/db"
	"booky-backend/internal/shared"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var (
	ErrInDatabase               = errors.New("database error")
	ErrOrderIsNotPending        = errors.New("order is not pending")
	ErrOrderDoesNotExist        = errors.New("order does not exist")
	ErrPaymentNotFound          = errors.New("payment not found")
	ErrInvalidPaymentTransition = errors.New("invalid payment transition")
	ErrAlreadyExists            = errors.New("payment already exists")
)

type OrderStatus string

const (
	OrderStatusPending OrderStatus = "pending"
	OrderStatusPaid    OrderStatus = "paid"
)

type PostgresRepository struct {
	db *db.DB
}

func NewPostgresRepo(db *db.DB) *PostgresRepository {
	return &PostgresRepository{
		db,
	}
}

// Create inserts a new payment (initial state: pending)
func (r *PostgresRepository) Create(ctx context.Context, req *CreatePaymentRequest) (*Payment, error) {
	// 1. Check idempotency first
	var existing Payment

	err := r.db.GetPool().QueryRow(ctx, `
		SELECT id, order_id, provider, provider_ref, status, created_at, updated_at
		FROM payments
		WHERE idempotency_key = $1
	`, req.IdempotencyKey).Scan(
		&existing.ID,
		&existing.OrderID,
		&existing.Provider,
		&existing.ProviderRef,
		&existing.Status,
		&existing.CreatedAt,
		&existing.UpdatedAt,
	)

	if err == nil {
		// already exists → return it (IMPORTANT: idempotent behavior)
		return &existing, nil
	}

	if err != pgx.ErrNoRows {
		shared.Log(shared.DEBUG, "failed to check idempotency: %v", err.Error())
		return nil, ErrInDatabase
	}

	// 2. Create new payment
	payment := &Payment{
		ID:             uuid.NewString(),
		OrderID:        req.OrderID,
		Provider:       "default", // or selected later
		ProviderRef:    "",
		IdempotencyKey: req.IdempotencyKey,
		Status:         StatusPending,
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
	}

	err = r.db.GetPool().QueryRow(ctx, `
		INSERT INTO payments (
			id,
			order_id,
			provider,
			provider_ref,
			idempotency_key,
			status,
			created_at,
			updated_at
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		RETURNING id
	`,
		payment.ID,
		payment.OrderID,
		payment.Provider,
		payment.ProviderRef,
		payment.IdempotencyKey,
		payment.Status,
		payment.CreatedAt,
		payment.UpdatedAt,
	).Scan(&payment.ID)

	if err != nil {
		shared.Log(shared.DEBUG, "failed to create payment: %v", err.Error())
		return nil, ErrInDatabase
	}

	return payment, nil
}

// GetByID fetches payment by internal ID
func (r *PostgresRepository) GetByID(ctx context.Context, id string) (*Payment, error) {
	return nil, nil
}

// GetByProviderRef fetches payment using external provider reference
func (r *PostgresRepository) GetByProviderRef(ctx context.Context, provider string, providerRef string) (*Payment, error) {
	return nil, nil
}

// UpdateStatus safely updates payment status (pending → paid/failed/etc.)
func (r *PostgresRepository) UpdateStatus(ctx context.Context, id string, newStatus PaymentStatus) error {
	return nil
}

// UpdateStatusByProviderRef is used by webhook (idempotent external update)
func (r *PostgresRepository) UpdateStatusByProviderRef(ctx context.Context, provider string, providerRef string, newStatus PaymentStatus) error {
	return nil
}

// LockForUpdate locks payment row for safe transitions
func (r *PostgresRepository) LockForUpdate(ctx context.Context, id string) (*Payment, error) {
	return nil, nil
}

// ListByOrderID (useful for UI / debugging)
func (r *PostgresRepository) ListByOrderID(ctx context.Context, orderID string) ([]Payment, error) {
	return nil, nil
}
