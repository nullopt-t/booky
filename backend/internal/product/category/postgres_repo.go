package category

import (
	"booky-backend/internal/model"
	"booky-backend/pkg/api"
	"booky-backend/pkg/database"
	"context"
)

type CategoryRepository interface {
	Create(ctx context.Context, qe database.QueryExecutor, name string) (*model.ProductCategory, error)
	GetAll(ctx context.Context, qe database.QueryExecutor, q *api.PageQuery) ([]*model.ProductCategory, *api.Page, error)
}

type Repository struct {
}

func NewPostgresRepository() *Repository {
	return &Repository{}
}

func (r *Repository) Create(ctx context.Context, qe database.QueryExecutor, name string) (*model.ProductCategory, error) {
	var category model.ProductCategory
	err := qe.QueryRow(ctx, "INSERT INTO categories (name) VALUES ($1) RETURNING id, name, created_at, updated_at", name).Scan(&category.ID, &category.Name, &category.CreatedAt, &category.UpdatedAt)
	if err != nil {
		return nil, database.MapError(err)
	}
	return &category, nil
}

func (r *Repository) GetAll(ctx context.Context, qe database.QueryExecutor, q *api.PageQuery) ([]*model.ProductCategory, *api.Page, error) {
	offset := (q.Page - 1) * q.Limit
	rows, err := qe.Query(ctx, `
		SELECT id, name, created_at, updated_at
		FROM categories
		ORDER BY name
		LIMIT $1 OFFSET $2
	`, q.Limit, offset)
	if err != nil {
		return nil, nil, database.MapError(err)
	}
	defer rows.Close()

	categories := make([]*model.ProductCategory, 0)
	for rows.Next() {
		var category model.ProductCategory
		err = rows.Scan(&category.ID, &category.Name, &category.CreatedAt, &category.UpdatedAt)
		if err != nil {
			return nil, nil, database.MapError(err)
		}
		categories = append(categories, &category)
	}

	var total int64
	err = qe.QueryRow(ctx, "SELECT COUNT(*) FROM categories").Scan(&total)
	if err != nil {
		return nil, nil, database.MapError(err)
	}

	page := &api.Page{
		Index: q.Page,
		Limit: q.Limit,
		Total: int(total),
	}

	return categories, page, nil
}
