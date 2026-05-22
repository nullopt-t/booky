package product

import (
	"booky-backend/internal/db"
	"booky-backend/internal/model"
	"booky-backend/internal/trans"
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type PostgresRepo struct {
}

func NewPostgresRepository() ProductRepository {
	return &PostgresRepo{}
}

func (r *PostgresRepo) Create(ctx context.Context, db db.DBQE, p *model.Product) (*model.Product, error) {
	var createdProduct model.Product
	err := db.QueryRow(ctx,
		`INSERT INTO products (title, price)
		 VALUES ($1, $2) RETURNING id, title, price, created_at, updated_at`,
		p.Title, p.Price,
	).Scan(&createdProduct.ID, &createdProduct.Title, &createdProduct.Price, &createdProduct.CreatedAt, &createdProduct.UpdatedAt)
	if err != nil {
		return nil, ErrInDatabase
	}
	return &createdProduct, nil
}

func (r *PostgresRepo) Save(ctx context.Context, db db.DBQE, p *model.Product) (*model.Product, error) {
	var updatedProduct model.Product
	err := db.QueryRow(ctx, "UPDATE products SET title = $1, price = $2 WHERE id = $3 RETURNING id, total, price, created_at, updated_at", p.Title, p.Price, p.ID).Scan(&updatedProduct.ID, &updatedProduct.Title, &updatedProduct.Price, &updatedProduct.CreatedAt, &updatedProduct.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrProductNotFount
		}
		return nil, ErrInDatabase
	}
	return &updatedProduct, nil
}

func (r *PostgresRepo) GetByID(ctx context.Context, db db.DBQE, productID uuid.UUID) (*model.Product, error) {
	var p model.Product

	err := db.QueryRow(ctx,
		`SELECT id, title, price, created_at, updated_at FROM products WHERE id=$1`,
		productID,
	).Scan(&p.ID, &p.Title, &p.Price, &p.CreatedAt, &p.UpdatedAt)

	if err != nil {
		return nil, ErrInDatabase
	}

	return &p, nil
}

func (r *PostgresRepo) GetAll(ctx context.Context, db db.DBQE, q trans.PaginationQuery) ([]*model.Product, *trans.Page, error) {
	offset := (q.Page - 1) * q.Limit
	rows, err := db.Query(ctx,
		`SELECT id, title, price, created_at, updated_at FROM products LIMIT $1 OFFSET $2`, q.Limit, offset)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var products = make([]*model.Product, 0, q.Limit)
	for rows.Next() {
		var p model.Product
		rows.Scan(&p.ID, &p.Title, &p.Price, &p.CreatedAt, &p.UpdatedAt)
		products = append(products, &p)
	}

	if err = rows.Err(); err != nil {
		return nil, nil, err
	}

	// query the products count
	var count int
	err = db.QueryRow(ctx, "SELECT COUNT(*) FROM products").Scan(&count)
	if err != nil {
		return nil, nil, err
	}

	resultPage := &trans.Page{
		Index: q.Page,
		Limit: q.Limit,
		Total: count,
	}

	return products, resultPage, nil
}
