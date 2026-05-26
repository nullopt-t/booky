package product

import (
	"booky-backend/internal/model"
	"booky-backend/pkg/api"
	"booky-backend/pkg/database"
	"context"

	"github.com/google/uuid"
)

type ProductRepository interface {
	Create(ctx context.Context, db database.QueryExecutor, p *model.Product) (*model.Product, error)
	Save(ctx context.Context, db database.QueryExecutor, p *model.Product) (*model.Product, error)
	GetByID(ctx context.Context, db database.QueryExecutor, productID uuid.UUID) (*model.Product, error)
	GetAll(ctx context.Context, db database.QueryExecutor, q api.PageQuery) ([]*model.Product, *api.Page, error)
}

type PostgresRepo struct {
}

func NewPostgresRepository() ProductRepository {
	return &PostgresRepo{}
}

func (r *PostgresRepo) Create(ctx context.Context, qe database.QueryExecutor, p *model.Product) (*model.Product, error) {
	var createdProduct model.Product
	err := qe.QueryRow(ctx,
		`INSERT INTO products (title, price)
		 VALUES ($1, $2) RETURNING id, title, price, created_at, updated_at`,
		p.Title, p.Price,
	).Scan(&createdProduct.ID, &createdProduct.Title, &createdProduct.Price, &createdProduct.CreatedAt, &createdProduct.UpdatedAt)
	if err != nil {
		return nil, database.MapError(err)
	}
	return &createdProduct, nil
}

func (r *PostgresRepo) Save(ctx context.Context, qe database.QueryExecutor, p *model.Product) (*model.Product, error) {
	var updatedProduct model.Product
	err := qe.QueryRow(ctx, "UPDATE products SET title = $1, price = $2 WHERE id = $3 RETURNING id, total, price, created_at, updated_at", p.Title, p.Price, p.ID).Scan(&updatedProduct.ID, &updatedProduct.Title, &updatedProduct.Price, &updatedProduct.CreatedAt, &updatedProduct.UpdatedAt)
	if err != nil {
		return nil, database.MapError(err)
	}
	return &updatedProduct, nil
}

func (r *PostgresRepo) GetByID(ctx context.Context, qe database.QueryExecutor, productID uuid.UUID) (*model.Product, error) {
	var p model.Product

	err := qe.QueryRow(ctx,
		`SELECT id, title, price, created_at, updated_at FROM products WHERE id=$1`,
		productID,
	).Scan(&p.ID, &p.Title, &p.Price, &p.CreatedAt, &p.UpdatedAt)

	if err != nil {
		return nil, database.MapError(err)
	}

	return &p, nil
}

func (r *PostgresRepo) GetAll(ctx context.Context, qe database.QueryExecutor, q api.PageQuery) ([]*model.Product, *api.Page, error) {
	offset := (q.Page - 1) * q.Limit
	rows, err := qe.Query(ctx,
		`SELECT id, title, price, created_at, updated_at FROM products LIMIT $1 OFFSET $2`, q.Limit, offset)
	if err != nil {
		return nil, nil, database.MapError(err)
	}
	defer rows.Close()

	var products = make([]*model.Product, 0, q.Limit)
	for rows.Next() {
		var p model.Product
		rows.Scan(&p.ID, &p.Title, &p.Price, &p.CreatedAt, &p.UpdatedAt)
		products = append(products, &p)
	}

	if err = rows.Err(); err != nil {
		return nil, nil, database.MapError(err)
	}

	// query the products count
	var count int
	err = qe.QueryRow(ctx, "SELECT COUNT(*) FROM products").Scan(&count)
	if err != nil {
		return nil, nil, database.MapError(err)
	}

	resultPage := &api.Page{
		Index: q.Page,
		Limit: q.Limit,
		Total: count,
	}

	return products, resultPage, nil
}
