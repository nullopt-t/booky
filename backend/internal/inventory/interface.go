package inventory

import (
	"booky-backend/pkg/database"
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type InventoryRepository interface {
	Reserve(ctx context.Context, qe database.DBQE, productID uuid.UUID, quantity int) error
	Release(ctx context.Context, qe database.DBQE, roductID uuid.UUID, quantity int) error
	GetAvailable(ctx context.Context, qe database.DBQE, productID uuid.UUID) (int, error)
	GetReserved(ctx context.Context, qe database.DBQE, productID uuid.UUID) (int, error)
}

type InventoryService interface {
	GetAvailable(ctx context.Context, productID uuid.UUID) (int, error)
	GetReserved(ctx context.Context, productID uuid.UUID) (int, error)
}

type InventoryHandler interface {
	GetAvailable(c *gin.Context)
	GetReserved(c *gin.Context)
}
