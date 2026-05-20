package inventory

import (
	"booky-backend/internal/db"
	"context"

	"github.com/google/uuid"
)

type InventoryRepository interface {
	Reserve(ctx context.Context, db db.Tx, productID uuid.UUID, quantity int) error
	Release(ctx context.Context, pdb db.Tx, roductID uuid.UUID, quantity int) error
	GetAvailable(ctx context.Context, db db.DBQE, productID uuid.UUID) (int, error)
	GetReserved(ctx context.Context, db db.DBQE, productID uuid.UUID) (int, error)
}
