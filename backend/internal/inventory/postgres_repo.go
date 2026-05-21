package inventory

import (
	"booky-backend/internal/db"
	"booky-backend/internal/shared"
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Repository struct {
}

func NewPostgresRepository() InventoryRepository {
	return &Repository{}
}

func (r *Repository) Reserve(ctx context.Context, db db.Tx, productID uuid.UUID, quantity int) error {
	var available_quantity int
	err := db.QueryRow(ctx, "SELECT available_quantity FROM inventories WHERE product_id = $1 FOR UPDATE", productID).Scan(&available_quantity)
	if err != nil {
		return ErrInDatabase
	}

	if available_quantity < quantity {
		return ErrInsufficientQuantity
	}

	// reserve the product
	_, err = db.Exec(ctx, "UPDATE inventories SET reserved_quantity += $1, available_quantity -= $1 WHERE product_id = $2", quantity, productID)
	if err != nil {
		shared.Log(shared.ERROR, err.Error())
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		return ErrInDatabase
	}

	return nil
}

func (r *Repository) Release(ctx context.Context, db db.Tx, productID uuid.UUID, quantity int) error {
	var reserved_quantity int
	err := db.QueryRow(ctx, "SELECT reserved_quantity FROM inventories WHERE product_id = $1 FOR UPDATE", productID).Scan(&reserved_quantity)
	if err != nil {
		shared.Log(shared.ERROR, err.Error())
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		return ErrInDatabase
	}

	if reserved_quantity < quantity {
		return ErrInsufficientQuantity
	}

	// reserve the product
	_, err = db.Exec(ctx, "UPDATE inventories SET reserved_quantity -= $1, available_quantity += $1 WHERE product_id = $2", quantity, productID)
	if err != nil {
		shared.Log(shared.ERROR, err.Error())
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		return ErrInDatabase
	}

	return nil
}

func (r *Repository) GetAvailable(ctx context.Context, db db.DBQE, productID uuid.UUID) (int, error) {
	var available_quantity int
	err := db.QueryRow(ctx, "SELECT available_quantity FROM inventories WHERE product_id = $1 ", productID).Scan(&available_quantity)
	if err != nil {
		shared.Log(shared.ERROR, err.Error())
		if errors.Is(err, pgx.ErrNoRows) {
			return available_quantity, ErrNotFound
		}
		return available_quantity, ErrInDatabase
	}
	return available_quantity, nil
}

func (r *Repository) GetReserved(ctx context.Context, db db.DBQE, productID uuid.UUID) (int, error) {
	var reserved_quantity int
	err := db.QueryRow(ctx, "SELECT reserved_quantity FROM inventories WHERE product_id = $1 ", productID).Scan(&reserved_quantity)
	if err != nil {
		shared.Log(shared.ERROR, err.Error())
		if errors.Is(err, pgx.ErrNoRows) {
			return reserved_quantity, ErrNotFound
		}
		return reserved_quantity, ErrInDatabase
	}
	return reserved_quantity, nil
}
