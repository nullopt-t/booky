package product

import (
	"booky-backend/internal/model"
	"booky-backend/pkg/api"
	"booky-backend/pkg/database"
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
)

// ── Mocks ─────────────────────────────────────────────────────────────────────

type MockProductRepository struct {
	CreateFn  func(ctx context.Context, db database.QueryExecutor, product *model.Product) (*model.Product, error)
	SaveFn    func(ctx context.Context, db database.QueryExecutor, product *model.Product) (*model.Product, error)
	GetByIDFn func(ctx context.Context, db database.QueryExecutor, id uuid.UUID) (*model.Product, error)
	GetAllFn  func(ctx context.Context, db database.QueryExecutor, q api.PageQuery) ([]*model.Product, *api.Page, error)
}

func (m *MockProductRepository) Create(ctx context.Context, db database.QueryExecutor, product *model.Product) (*model.Product, error) {
	if m.CreateFn == nil {
		panic("CreateFn is not set")
	}
	return m.CreateFn(ctx, db, product)
}
func (m *MockProductRepository) Save(ctx context.Context, db database.QueryExecutor, product *model.Product) (*model.Product, error) {
	if m.SaveFn == nil {
		panic("SaveFn is not set")
	}
	return m.SaveFn(ctx, db, product)
}
func (m *MockProductRepository) GetByID(ctx context.Context, db database.QueryExecutor, id uuid.UUID) (*model.Product, error) {
	if m.GetByIDFn == nil {
		panic("GetByIDFn is not set")
	}
	return m.GetByIDFn(ctx, db, id)
}
func (m *MockProductRepository) GetAll(ctx context.Context, db database.QueryExecutor, q api.PageQuery) ([]*model.Product, *api.Page, error) {
	if m.GetAllFn == nil {
		panic("GetAllFn is not set")
	}
	return m.GetAllFn(ctx, db, q)
}

type MockInventoryRepository struct{}

type MockRunner struct {
	WithTxFn func(ctx context.Context, fn func(tx database.QueryExecutor) error) error
	WithDBFn func(ctx context.Context, fn func(pool database.QueryExecutor) error) error
}

func (m *MockRunner) WithTx(ctx context.Context, fn func(tx database.QueryExecutor) error) error {
	if m.WithTxFn == nil {
		panic("WithTxFn is not set")
	}
	return m.WithTxFn(ctx, fn)
}
func (m *MockRunner) WithDB(ctx context.Context, fn func(pool database.QueryExecutor) error) error {
	if m.WithDBFn == nil {
		panic("WithDBFn is not set")
	}
	return m.WithDBFn(ctx, fn)
}

// ── Helpers ───────────────────────────────────────────────────────────────────

func execTx(_ context.Context, fn func(database.QueryExecutor) error) error { return fn(nil) }
func noDB() database.QueryExecutor                                          { return nil }

// ── TestCreate ────────────────────────────────────────────────────────────────

func TestCreate(t *testing.T) {
	runner := &MockRunner{WithTxFn: execTx}

	t.Run("success: returns created product", func(t *testing.T) {
		req := CreateProductRequest{Title: "Book", Price: 99}
		repo := &MockProductRepository{
			CreateFn: func(_ context.Context, _ database.QueryExecutor, p *model.Product) (*model.Product, error) {
				if p.Title != req.Title || p.Price != req.Price {
					t.Fatalf("unexpected product passed to Create: %+v", p)
				}
				return &model.Product{ID: uuid.New(), Title: p.Title, Price: p.Price}, nil
			},
		}

		product, err := NewService(runner, repo, nil).Create(context.Background(), req)
		if err != nil {
			t.Fatal("unexpected error:", err)
		}
		if product.Title != req.Title || product.Price != req.Price {
			t.Fatalf("returned product does not match request: %+v", product)
		}
	})

	t.Run("repo error: returns error", func(t *testing.T) {
		repo := &MockProductRepository{
			CreateFn: func(_ context.Context, _ database.QueryExecutor, _ *model.Product) (*model.Product, error) {
				return nil, fmt.Errorf("db error")
			},
		}

		if _, err := NewService(runner, repo, nil).Create(context.Background(), CreateProductRequest{Title: "Book", Price: 99}); err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

// ── TestUpdate ────────────────────────────────────────────────────────────────

func TestUpdate(t *testing.T) {
	runner := &MockRunner{WithTxFn: execTx, WithDBFn: func(_ context.Context, fn func(pool database.QueryExecutor) error) error { return fn(nil) }}

	existingProduct := &model.Product{ID: uuid.New(), Title: "Old Title", Price: 50}

	getByID := func(_ context.Context, _ database.QueryExecutor, _ uuid.UUID) (*model.Product, error) {
		return existingProduct, nil
	}

	t.Run("updates both title and price", func(t *testing.T) {
		newTitle, newPrice := "New Title", 200
		repo := &MockProductRepository{
			GetByIDFn: getByID,
			SaveFn: func(_ context.Context, _ database.QueryExecutor, p *model.Product) (*model.Product, error) {
				if p.Title != newTitle {
					t.Fatalf("expected title %q, got %q", newTitle, p.Title)
				}
				if p.Price != newPrice {
					t.Fatalf("expected price %d, got %d", newPrice, p.Price)
				}
				return p, nil
			},
		}

		if _, err := NewService(runner, repo, nil).Update(context.Background(), existingProduct.ID, UpdateProductRequest{
			Title: &newTitle, Price: &newPrice,
		}); err != nil {
			t.Fatal("unexpected error:", err)
		}
	})

	t.Run("partial update: only price", func(t *testing.T) {
		newPrice := 300
		repo := &MockProductRepository{
			GetByIDFn: getByID,
			SaveFn: func(_ context.Context, _ database.QueryExecutor, p *model.Product) (*model.Product, error) {
				if p.Title != existingProduct.Title {
					t.Fatalf("title should be unchanged, got %q", p.Title)
				}
				if p.Price != newPrice {
					t.Fatalf("expected price %d, got %d", newPrice, p.Price)
				}
				return p, nil
			},
		}

		if _, err := NewService(runner, repo, nil).Update(context.Background(), existingProduct.ID, UpdateProductRequest{
			Price: &newPrice,
		}); err != nil {
			t.Fatal("unexpected error:", err)
		}
	})

	t.Run("product not found: returns error", func(t *testing.T) {
		repo := &MockProductRepository{
			GetByIDFn: func(_ context.Context, _ database.QueryExecutor, _ uuid.UUID) (*model.Product, error) {
				return nil, fmt.Errorf("not found")
			},
		}

		if _, err := NewService(runner, repo, nil).Update(context.Background(), uuid.New(), UpdateProductRequest{}); err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("save fails: returns error", func(t *testing.T) {
		repo := &MockProductRepository{
			GetByIDFn: getByID,
			SaveFn: func(_ context.Context, _ database.QueryExecutor, _ *model.Product) (*model.Product, error) {
				return nil, fmt.Errorf("save failed")
			},
		}

		newTitle := "anything"
		if _, err := NewService(runner, repo, nil).Update(context.Background(), existingProduct.ID, UpdateProductRequest{
			Title: &newTitle,
		}); err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

// ── TestGetAll ────────────────────────────────────────────────────────────────

func TestGetAll(t *testing.T) {
	runner := &MockRunner{WithDBFn: func(_ context.Context, fn func(pool database.QueryExecutor) error) error { return fn(nil) }}

	t.Run("success: returns products and page", func(t *testing.T) {
		products := []*model.Product{{ID: uuid.New()}, {ID: uuid.New()}}
		page := &api.Page{Total: 2}
		repo := &MockProductRepository{
			GetAllFn: func(_ context.Context, _ database.QueryExecutor, _ api.PageQuery) ([]*model.Product, *api.Page, error) {
				return products, page, nil
			},
		}

		got, gotPage, err := NewService(runner, repo, nil).GetAll(context.Background(), api.PageQuery{})
		if err != nil {
			t.Fatal("unexpected error:", err)
		}
		if len(got) != 2 || gotPage.Total != 2 {
			t.Fatalf("unexpected result: %+v, %+v", got, gotPage)
		}
	})

	t.Run("repo error: returns error", func(t *testing.T) {
		repo := &MockProductRepository{
			GetAllFn: func(_ context.Context, _ database.QueryExecutor, _ api.PageQuery) ([]*model.Product, *api.Page, error) {
				return nil, nil, fmt.Errorf("db error")
			},
		}

		if _, _, err := NewService(runner, repo, nil).GetAll(context.Background(), api.PageQuery{}); err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

// ── TestGetByID ───────────────────────────────────────────────────────────────

func TestGetByID(t *testing.T) {
	runner := &MockRunner{WithDBFn: func(_ context.Context, fn func(pool database.QueryExecutor) error) error { return fn(nil) }}

	t.Run("success: returns product", func(t *testing.T) {
		productID := uuid.New()
		repo := &MockProductRepository{
			GetByIDFn: func(_ context.Context, _ database.QueryExecutor, id uuid.UUID) (*model.Product, error) {
				return &model.Product{ID: id}, nil
			},
		}

		product, err := NewService(runner, repo, nil).GetByID(context.Background(), productID)
		if err != nil {
			t.Fatal("unexpected error:", err)
		}
		if product.ID != productID {
			t.Fatalf("expected ID %v, got %v", productID, product.ID)
		}
	})

	t.Run("not found: returns error", func(t *testing.T) {
		repo := &MockProductRepository{
			GetByIDFn: func(_ context.Context, _ database.QueryExecutor, _ uuid.UUID) (*model.Product, error) {
				return nil, fmt.Errorf("not found")
			},
		}

		if _, err := NewService(runner, repo, nil).GetByID(context.Background(), uuid.New()); err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}
