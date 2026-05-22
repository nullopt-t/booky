package inventory

import (
	"booky-backend/pkg/database"
	"context"

	"github.com/google/uuid"
)

type InventoryRepository interface {
	Reserve(ctx context.Context, qe database.DBQE, productID uuid.UUID, quantity int) error
	Release(ctx context.Context, qe database.DBQE, roductID uuid.UUID, quantity int) error
	GetAvailable(ctx context.Context, qe database.DBQE, productID uuid.UUID) (int, error)
	GetReserved(ctx context.Context, qe database.DBQE, productID uuid.UUID) (int, error)
}
