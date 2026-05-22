package inventory

import (
	"booky-backend/pkg/database"
	"context"

	"github.com/google/uuid"
)

type Repository struct {
}

func NewPostgresRepository() InventoryRepository {
	return &Repository{}
}

func (r *Repository) Reserve(ctx context.Context, qe database.QueryExecutor, productID uuid.UUID, quantity int) error {
	var available_quantity int
	err := qe.QueryRow(ctx, "SELECT available_quantity FROM inventories WHERE product_id = $1 FOR UPDATE", productID).Scan(&available_quantity)
	if err != nil {
		return database.MapError(err)
	}

	// reserve the product
	_, err = qe.Exec(ctx, "UPDATE inventories SET reserved_quantity += $1, available_quantity -= $1 WHERE product_id = $2 AND available_quantity >= $1", quantity, productID)
	if err != nil {
		return database.MapError(err)
	}

	return nil
}

func (r *Repository) Release(ctx context.Context, qe database.QueryExecutor, productID uuid.UUID, quantity int) error {
	var reserved_quantity int
	err := qe.QueryRow(ctx, "SELECT reserved_quantity FROM inventories WHERE product_id = $1 FOR UPDATE", productID).Scan(&reserved_quantity)
	if err != nil {
		return database.MapError(err)
	}

	// reserve the product
	_, err = qe.Exec(ctx, "UPDATE inventories SET reserved_quantity -= $1, available_quantity += $1 WHERE product_id = $2 AND reserved_quantity >= $1", quantity, productID)
	if err != nil {
		return database.MapError(err)
	}

	return nil
}

func (r *Repository) GetAvailable(ctx context.Context, qe database.QueryExecutor, productID uuid.UUID) (int, error) {
	var available_quantity int
	err := qe.QueryRow(ctx, "SELECT available_quantity FROM inventories WHERE product_id = $1 ", productID).Scan(&available_quantity)
	if err != nil {
		return available_quantity, database.MapError(err)
	}
	return available_quantity, nil
}

func (r *Repository) GetReserved(ctx context.Context, qe database.QueryExecutor, productID uuid.UUID) (int, error) {
	var reserved_quantity int
	err := qe.QueryRow(ctx, "SELECT reserved_quantity FROM inventories WHERE product_id = $1 ", productID).Scan(&reserved_quantity)
	if err != nil {
		return reserved_quantity, database.MapError(err)
	}
	return reserved_quantity, nil
}
