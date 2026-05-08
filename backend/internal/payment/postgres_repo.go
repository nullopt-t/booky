package payment

import (
	"booky-backend/internal/db"
	"booky-backend/internal/order"
	"booky-backend/internal/shared"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

var (
	ErrInDatabase               = errors.New("database error")
	ErrOrderIsNotPending        = errors.New("order is not pending")
	ErrOrderDoesNotExist        = errors.New("order does not exist")
	ErrPaymentNotFound          = errors.New("payment not found")
	ErrInvalidPaymentTransition = errors.New("invalid payment transition")
)

type PostgresRepository struct {
	db *db.DB
}

func NewPostgresRepo(db *db.DB) *PostgresRepository {
	return &PostgresRepository{
		db,
	}
}

func (r *PostgresRepository) GetByID(ctx context.Context, id string) (*Payment, error) {
	var p Payment
	err := r.db.GetPool().QueryRow(ctx, `
	SELECT id, order_id, amount, status, provider, provider_ref, created_at, updated_at
	FROM payments
	WHERE id = $1 FOR UPDATE`, id).Scan(&p.ID, &p.OrderID, &p.Amount, &p.Status, &p.Provider, &p.ProviderRef, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrPaymentNotFound
		}
		shared.Log(shared.DEBUG, "failed to get payment: %v", err.Error())
		return nil, ErrInDatabase
	}
	return &p, nil
}

func (r *PostgresRepository) Create(ctx context.Context, orderID string, provider string, providerRef string) (*Payment, error) {
	tx, err := r.db.GetPool().Begin(ctx)
	if err != nil {
		shared.Log(shared.ERROR, "failed to begin transaction: %v", err)
		return nil, err
	}
	defer tx.Rollback(ctx)

	type OrderWithStatusAndTotal struct {
		Status     order.OrderStatus
		TotalPrice int
		OrderID    string
	}

	var orderDetails OrderWithStatusAndTotal
	err = tx.QueryRow(ctx, "SELECT id, status, total_price FROM orders WHERE id = $1", orderID).Scan(&orderDetails.OrderID, &orderDetails.Status, &orderDetails.TotalPrice)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			shared.Log(shared.ERROR, "order does not exist")
			return nil, ErrOrderDoesNotExist
		}
		shared.Log(shared.ERROR, "failed to get order status: %v", err)
		return nil, ErrInDatabase
	}

	if orderDetails.Status != order.OrderStatusPending {
		return nil, ErrOrderIsNotPending
	}

	var payment Payment
	err = tx.QueryRow(ctx, `
    INSERT INTO payments (order_id, amount, status, provider, provider_ref) 
    VALUES ($1, $2, $3, $4, $5) 
    RETURNING id, order_id, amount, status, provider, provider_ref, created_at, updated_at`,
		orderID, orderDetails.TotalPrice, string(PaymentStatusPending), provider, providerRef,
	).Scan(&payment.ID, &payment.OrderID, &payment.Amount, &payment.Status, &payment.Provider, &payment.ProviderRef, &payment.CreatedAt, &payment.UpdatedAt)
	if err != nil {
		shared.Log(shared.ERROR, "failed to create payment: %v", err)
		return nil, ErrInDatabase
	}

	err = tx.Commit(ctx)
	if err != nil {
		shared.Log(shared.ERROR, "failed to commit transaction: %v", err)
		return nil, ErrInDatabase
	}

	return &payment, nil
}

func (r *PostgresRepository) GetByProviderRef(ctx context.Context, provider string, providerRef string) (*Payment, error) {
	var p Payment
	err := r.db.GetPool().QueryRow(ctx, `
	SELECT id, order_id, amount, status, provider, provider_ref, created_at, updated_at
	FROM payments
	WHERE provider = $1 AND provider_ref = $2`, provider, providerRef).Scan(&p.ID, &p.OrderID, &p.Amount, &p.Status, &p.Provider, &p.ProviderRef, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrPaymentNotFound
		}
		shared.Log(shared.DEBUG, "failed to get payment: %v", err.Error())
		return nil, ErrInDatabase
	}
	return &p, nil
}

func (r *PostgresRepository) TransitionStatus(ctx context.Context, id string, oldStatus PaymentStatus, newStatus PaymentStatus) error {
	result, err := r.db.GetPool().Exec(ctx, `
	UPDATE payments
	SET status = $1, updated_at = NOW()
	WHERE id = $2 AND status = $3`, newStatus, id, oldStatus)
	if err != nil {
		shared.Log(shared.DEBUG, "failed to update payment status: %v", err.Error())
		return ErrInDatabase
	}
	if rowsAffected := result.RowsAffected(); rowsAffected == 0 {
		return ErrInvalidPaymentTransition
	}
	return nil
}

func (r *PostgresRepository) TransitionStatusByProviderRef(ctx context.Context, provider string, providerRef string, newStatus PaymentStatus) error {
	result, err := r.db.GetPool().Exec(ctx, `
	UPDATE payments
	SET status = $1, updated_at = NOW()
	WHERE provider = $2 AND provider_ref = $3`, newStatus, provider, providerRef)
	if err != nil {
		shared.Log(shared.DEBUG, "failed to update payment status: %v", err.Error())
		return ErrInDatabase
	}
	if rowsAffected := result.RowsAffected(); rowsAffected == 0 {
		return ErrInvalidPaymentTransition
	}
	return nil
}
