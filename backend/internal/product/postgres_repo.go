package product

import (
	"booky-backend/internal/db"
	"booky-backend/internal/domain"
	"booky-backend/internal/utils"
	"context"
	"fmt"
)

type PostgresRepo struct {
	db *db.DB
}

func NewPostgresRepo(db *db.DB) *PostgresRepo {
	return &PostgresRepo{
		db,
	}
}

func (r *PostgresRepo) Create(ctx context.Context, req CreateProductRequest) (*domain.Product, error) {
	var p domain.Product
	err := r.db.GetPool().QueryRow(ctx,
		`INSERT INTO products (title, price, stock)
		 VALUES ($1, $2, $3) RETURNING id, title, price, stock, created_at, updated_at`,
		req.Title, req.Price, req.Stock,
	).Scan(&p.ID, &p.Title, &p.Price, &p.Stock, &p.CreatedAt, &p.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("insert product: %w", err)
	}
	return &p, nil
}

func (r *PostgresRepo) Update(ctx context.Context, id string, req UpdateProductRequest) (*domain.Product, error) {
	var p domain.Product
	query := `UPDATE products SET title=COALESCE($2, title), price=COALESCE($3, price), stock=COALESCE($4, stock), updated_at=now() WHERE id=$1 RETURNING id, title, price, stock, created_at, updated_at`
	args := []interface{}{id, req.Title, req.Price, req.Stock}

	err := r.db.GetPool().QueryRow(ctx, query, args...).Scan(&p.ID, &p.Title, &p.Price, &p.Stock, &p.CreatedAt, &p.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("update product: %w", err)
	}

	return &p, nil
}

func (r *PostgresRepo) GetByID(ctx context.Context, id string) (*domain.Product, error) {
	var p domain.Product

	err := r.db.GetPool().QueryRow(ctx,
		`SELECT id, title, price, stock, created_at, updated_at FROM products WHERE id=$1`,
		id,
	).Scan(&p.ID, &p.Title, &p.Price, &p.Stock, &p.CreatedAt, &p.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("get product by id : %w", err)
	}

	return &p, nil
}

func (r *PostgresRepo) GetAll(ctx context.Context, q utils.PaginationQuery) (*utils.PageResult[domain.Product], error) {
	offset := (q.Page - 1) * q.Limit
	rows, err := r.db.GetPool().Query(ctx,
		`SELECT id, title, price, stock, created_at, updated_at FROM products LIMIT $1 OFFSET $2`, q.Limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products = []domain.Product{}
	for rows.Next() {
		var p domain.Product
		rows.Scan(&p.ID, &p.Title, &p.Price, &p.Stock, &p.CreatedAt, &p.UpdatedAt)
		products = append(products, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// query the products count
	var count int
	err = r.db.GetPool().QueryRow(ctx, "SELECT COUNT(*) FROM products").Scan(&count)
	if err != nil {
		return nil, err
	}

	return &utils.PageResult[domain.Product]{
		Items: products,
		Page:  q.Page,
		Limit: q.Limit,
		Total: count,
	}, nil
}

