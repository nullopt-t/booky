package product

import (
	"booky-backend/internal/model"
	"booky-backend/pkg/api"
	"booky-backend/pkg/database"
	"context"

	"github.com/google/uuid"
)

type InventoryRepository interface {
	GetAvailable(ctx context.Context, db database.QueryExecutor, productID uuid.UUID) (int, error)
	GetReserved(ctx context.Context, db database.QueryExecutor, productID uuid.UUID) (int, error)
}

type Service struct {
	dbExecuter    database.Runner
	productRepo   ProductRepository
	inventoryRepo InventoryRepository
}

func NewService(dbExecuter database.Runner, productRepo ProductRepository, inventoryRepo InventoryRepository) ProudctService {
	return &Service{dbExecuter, productRepo, inventoryRepo}
}

func (s *Service) Create(ctx context.Context, req CreateProductRequest) (*model.Product, error) {
	var createdProduct *model.Product
	err := s.dbExecuter.WithTx(ctx, func(tx database.QueryExecutor) error {
		var err error
		createdProduct, err = s.productRepo.Create(ctx, tx, &model.Product{
			Title: req.Title,
			Price: req.Price,
		})
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return createdProduct, nil
}

func (s *Service) Update(ctx context.Context, productID uuid.UUID, req UpdateProductRequest) (*model.Product, error) {
	var existedProduct *model.Product
	err := s.dbExecuter.WithDB(ctx, func(db database.QueryExecutor) error {
		var err error
		existedProduct, err = s.productRepo.GetByID(ctx, db, productID)
		return err
	})
	if err != nil {
		return nil, err
	}

	if req.Price != nil {
		existedProduct.Price = *req.Price
	}

	if req.Title != nil {
		existedProduct.Title = *req.Title
	}

	var savedProduct *model.Product
	err = s.dbExecuter.WithTx(ctx, func(tx database.QueryExecutor) error {
		var err error
		savedProduct, err = s.productRepo.Save(ctx, tx, existedProduct)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return savedProduct, nil
}

func (s *Service) GetAll(ctx context.Context, q api.PageQuery) ([]*model.Product, *api.Page, error) {
	var products []*model.Product
	var page *api.Page
	err := s.dbExecuter.WithDB(ctx, func(db database.QueryExecutor) error {
		var err error
		products, page, err = s.productRepo.GetAll(ctx, db, q)
		return err
	})
	if err != nil {
		return nil, nil, err
	}
	return products, page, nil
}

func (s *Service) GetByID(ctx context.Context, productID uuid.UUID) (*model.Product, error) {
	var product *model.Product
	err := s.dbExecuter.WithDB(ctx, func(db database.QueryExecutor) error {
		var err error
		product, err = s.productRepo.GetByID(ctx, db, productID)
		return err
	})
	if err != nil {
		return nil, err
	}
	return product, nil
}
