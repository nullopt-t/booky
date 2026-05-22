package product

import (
	"booky-backend/internal/model"
	"booky-backend/pkg/api"
	"booky-backend/pkg/database"
	"context"

	"github.com/google/uuid"
)

type Service struct {
	tx            database.Runner
	productRepo   ProductRepository
	inventoryRepo InventoryRepository
}

func NewService(tx database.Runner, productRepo ProductRepository, inventoryRepo InventoryRepository) ProudctService {
	return &Service{tx, productRepo, inventoryRepo}
}

func (s *Service) Create(ctx context.Context, req CreateProductRequest) (*model.Product, error) {
	var createdProduct *model.Product
	err := s.tx.WithTx(ctx, func(tx database.QueryExecutor) error {
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
	existedProduct, err := s.productRepo.GetByID(ctx, s.tx.DB(), productID)
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
	err = s.tx.WithTx(ctx, func(tx database.QueryExecutor) error {
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
	return s.productRepo.GetAll(ctx, s.tx.DB(), q)
}

func (s *Service) GetByID(ctx context.Context, productID uuid.UUID) (*model.Product, error) {
	return s.productRepo.GetByID(ctx, s.tx.DB(), productID)
}
