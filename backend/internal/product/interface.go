package product

import (
	"booky-backend/internal/db"
	"booky-backend/internal/model"
	"booky-backend/internal/trans"
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type InventoryRepository interface {
	GetAvailable(ctx context.Context, db db.DBQE, productID uuid.UUID) (int, error)
	GetReserved(ctx context.Context, db db.DBQE, productID uuid.UUID) (int, error)
}

type ProductRepository interface {
	Create(ctx context.Context, db db.DBQE, p *model.Product) (*model.Product, error)
	Save(ctx context.Context, db db.DBQE, p *model.Product) (*model.Product, error)
	GetByID(ctx context.Context, db db.DBQE, productID uuid.UUID) (*model.Product, error)
	GetAll(ctx context.Context, db db.DBQE, q trans.PaginationQuery) ([]*model.Product, *trans.Page, error)
}

type ProudctService interface {
	Create(ctx context.Context, req CreateProductRequest) (*model.Product, error)
	Update(ctx context.Context, productID uuid.UUID, req UpdateProductRequest) (*model.Product, error)
	GetByID(ctx context.Context, productID uuid.UUID) (*model.Product, error)
	GetAll(ctx context.Context, q trans.PaginationQuery) ([]*model.Product, *trans.Page, error)
}

type ProductHandler interface {
	CreateProduct(c *gin.Context)
	GetProductByID(c *gin.Context)
	GetAllProducts(c *gin.Context)
}
