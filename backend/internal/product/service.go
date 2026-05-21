package product

import (
	"booky-backend/internal/db"
	"booky-backend/internal/model"
	"booky-backend/internal/trans"
	"context"

	"github.com/google/uuid"
)

type Service struct {
	tx            db.Runner
	productRepo   ProductRepository
	inventoryRepo InventoryRepository
}

func NewService(tx db.Runner, productRepo ProductRepository, inventoryRepo InventoryRepository) ProudctService {
	return &Service{tx, productRepo, inventoryRepo}
}

func (s *Service) Create(ctx context.Context, req CreateProductRequest) (*model.Product, error) {
	var createdProduct *model.Product
	err := s.tx.WithTx(ctx, func(tx db.DBQE) error {
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
	err = s.tx.WithTx(ctx, func(tx db.DBQE) error {
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

func (s *Service) GetAll(ctx context.Context, q trans.PaginationQuery) ([]*model.Product, *trans.Page, error) {
	return s.productRepo.GetAll(ctx, s.tx.DB(), q)
}

func (s *Service) GetByID(ctx context.Context, productID uuid.UUID) (*model.Product, error) {
	return s.productRepo.GetByID(ctx, s.tx.DB(), productID)
}
