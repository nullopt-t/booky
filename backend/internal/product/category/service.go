package category

import (
	"booky-backend/internal/model"
	"booky-backend/pkg/api"
	"booky-backend/pkg/database"
	"context"
)

type Service struct {
	runner database.Runner
	repo   CategoryRepository
}

func NewService(r database.Runner, repo CategoryRepository) *Service {
	return &Service{
		runner: r,
		repo:   repo,
	}
}

func (s *Service) Create(ctx context.Context, name string) (*model.ProductCategory, error) {
	var category *model.ProductCategory
	err := s.runner.WithDB(ctx, func(db database.QueryExecutor) error {
		var err error
		category, err = s.repo.Create(ctx, db, name)
		return err
	})
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (s *Service) GetAll(ctx context.Context, q *api.PageQuery) ([]*model.ProductCategory, *api.Page, error) {
	var categories []*model.ProductCategory
	var page *api.Page
	err := s.runner.WithDB(ctx, func(db database.QueryExecutor) error {
		var err error
		categories, page, err = s.repo.GetAll(ctx, db, q)
		return err
	})
	if err != nil {
		return nil, nil, err
	}
	return categories, page, nil
}
