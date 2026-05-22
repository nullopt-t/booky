package cart

import (
	"booky-backend/internal/db"
	"booky-backend/internal/model"
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
)

type MockCartRepository struct {
	CreateFn      func(ctx context.Context, db db.DBQE, userID uuid.UUID) (*model.Cart, error)
	GetByUserIDFn func(ctx context.Context, db db.DBQE, userID uuid.UUID) (*model.Cart, error)
	EmptyFn       func(ctx context.Context, db db.DBQE, userID uuid.UUID) error
	SaveFn        func(ctx context.Context, db db.DBQE, cart *model.Cart) error
}

func (m *MockCartRepository) Create(ctx context.Context, db db.DBQE, userID uuid.UUID) (*model.Cart, error) {
	if m.CreateFn == nil {
		panic("CreateFn is not set")
	}
	return m.CreateFn(ctx, db, userID)
}
func (m *MockCartRepository) GetByUserID(ctx context.Context, db db.DBQE, userID uuid.UUID) (*model.Cart, error) {
	if m.GetByUserIDFn == nil {
		panic("GetByUserIDFn is not set")
	}
	return m.GetByUserIDFn(ctx, db, userID)
}
func (m *MockCartRepository) Empty(ctx context.Context, db db.DBQE, userID uuid.UUID) error {
	if m.EmptyFn == nil {
		panic("EmptyFn is not set")
	}
	return m.EmptyFn(ctx, db, userID)
}
func (m *MockCartRepository) Save(ctx context.Context, db db.DBQE, cart *model.Cart) error {
	if m.SaveFn == nil {
		panic("SaveFn is not set")
	}
	return m.SaveFn(ctx, db, cart)
}

type MockRunner struct {
	WithTxFn func(ctx context.Context, fn func(tx db.DBQE) error) error
	DBFn     func() db.DBQE
}

func (m *MockRunner) WithTx(ctx context.Context, fn func(tx db.DBQE) error) error {
	if m.WithTxFn == nil {
		panic("WithTxFn is not set")
	}
	return m.WithTxFn(ctx, fn)
}
func (m *MockRunner) DB() db.DBQE {
	if m.DBFn == nil {
		panic("DBFn is not set")
	}
	return m.DBFn()
}

type MockProductRepository struct {
	GetByIDFn func(ctx context.Context, db db.DBQE, id uuid.UUID) (*model.Product, error)
}

func (m *MockProductRepository) GetByID(ctx context.Context, db db.DBQE, id uuid.UUID) (*model.Product, error) {
	if m.GetByIDFn == nil {
		panic("GetByIDFn is not set")
	}
	return m.GetByIDFn(ctx, db, id)
}

func execTx(ctx context.Context, fn func(tx db.DBQE) error) error { return fn(nil) }

func TestGetCart(t *testing.T) {
	runner := &MockRunner{DBFn: func() db.DBQE { return nil }}

	t.Run("success: returns items and correct total", func(t *testing.T) {
		p1, p2 := uuid.New(), uuid.New()
		repo := &MockCartRepository{
			GetByUserIDFn: func(_ context.Context, _ db.DBQE, _ uuid.UUID) (*model.Cart, error) {
				return &model.Cart{Items: []model.CartItem{
					{ProductID: p1, Quantity: 1},
					{ProductID: p2, Quantity: 2},
				}}, nil
			},
		}
		prices := map[uuid.UUID]int{p1: 100, p2: 50}
		productRepo := &MockProductRepository{
			GetByIDFn: func(_ context.Context, _ db.DBQE, id uuid.UUID) (*model.Product, error) {
				return &model.Product{ID: id, Price: prices[id]}, nil
			},
		}

		_, total, err := NewService(runner, repo, productRepo).GetCart(context.Background(), uuid.New())
		if err != nil {
			t.Fatal("unexpected error:", err)
		}
		// 1×100 + 2×50 = 200
		if total != 200 {
			t.Fatalf("expected total 200, got %d", total)
		}
	})

	t.Run("get fails, create succeeds: returns empty cart", func(t *testing.T) {
		userID := uuid.New()
		repo := &MockCartRepository{
			GetByUserIDFn: func(_ context.Context, _ db.DBQE, _ uuid.UUID) (*model.Cart, error) {
				return nil, db.ErrNotFound
			},
			CreateFn: func(_ context.Context, _ db.DBQE, id uuid.UUID) (*model.Cart, error) {
				return &model.Cart{UserID: id, Items: []model.CartItem{}}, nil
			},
		}

		cart, total, err := NewService(runner, repo, nil).GetCart(context.Background(), userID)
		if err != nil {
			t.Fatal("unexpected error:", err)
		}
		if cart.UserID != userID || total != 0 {
			t.Fatalf("expected empty cart for user %v, got %+v total %d", userID, cart, total)
		}
	})

	t.Run("get and create both fail: returns error", func(t *testing.T) {
		repo := &MockCartRepository{
			GetByUserIDFn: func(_ context.Context, _ db.DBQE, _ uuid.UUID) (*model.Cart, error) {
				return nil, fmt.Errorf("not found")
			},
			CreateFn: func(_ context.Context, _ db.DBQE, _ uuid.UUID) (*model.Cart, error) {
				return nil, fmt.Errorf("create failed")
			},
		}

		_, _, err := NewService(runner, repo, nil).GetCart(context.Background(), uuid.New())
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("product fetch fails: returns error", func(t *testing.T) {
		repo := &MockCartRepository{
			GetByUserIDFn: func(_ context.Context, _ db.DBQE, _ uuid.UUID) (*model.Cart, error) {
				return &model.Cart{Items: []model.CartItem{{ProductID: uuid.New(), Quantity: 1}}}, nil
			},
		}
		productRepo := &MockProductRepository{
			GetByIDFn: func(_ context.Context, _ db.DBQE, _ uuid.UUID) (*model.Product, error) {
				return nil, fmt.Errorf("product not found")
			},
		}

		_, _, err := NewService(runner, repo, productRepo).GetCart(context.Background(), uuid.New())
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestEmptyCart(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		runner := &MockRunner{WithTxFn: execTx}
		repo := &MockCartRepository{EmptyFn: func(_ context.Context, _ db.DBQE, _ uuid.UUID) error { return nil }}

		if err := NewService(runner, repo, nil).EmptyCart(context.Background(), uuid.New()); err != nil {
			t.Fatal("unexpected error:", err)
		}
	})

	t.Run("transaction failure", func(t *testing.T) {
		runner := &MockRunner{
			WithTxFn: func(_ context.Context, _ func(db.DBQE) error) error {
				return fmt.Errorf("tx failed")
			},
		}

		if err := NewService(runner, nil, nil).EmptyCart(context.Background(), uuid.New()); err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestAddItem(t *testing.T) {
	runner := &MockRunner{WithTxFn: execTx}

	t.Run("new product appended to cart", func(t *testing.T) {
		existing := uuid.New()
		repo := &MockCartRepository{
			GetByUserIDFn: func(_ context.Context, _ db.DBQE, _ uuid.UUID) (*model.Cart, error) {
				return &model.Cart{Items: []model.CartItem{{ProductID: existing, Quantity: 1}}}, nil
			},
			SaveFn: func(_ context.Context, _ db.DBQE, cart *model.Cart) error {
				if len(cart.Items) != 2 {
					t.Fatalf("expected 2 items, got %d", len(cart.Items))
				}
				return nil
			},
		}

		if _, err := NewService(runner, repo, nil).AddItem(context.Background(), uuid.New(), AddCartItemRequest{
			ProductID: uuid.New(), Quantity: 1,
		}); err != nil {
			t.Fatal("unexpected error:", err)
		}
	})

	t.Run("existing product quantity incremented", func(t *testing.T) {
		productID := uuid.New()
		repo := &MockCartRepository{
			GetByUserIDFn: func(_ context.Context, _ db.DBQE, _ uuid.UUID) (*model.Cart, error) {
				return &model.Cart{Items: []model.CartItem{{ProductID: productID, Quantity: 1}}}, nil
			},
			SaveFn: func(_ context.Context, _ db.DBQE, cart *model.Cart) error {
				if len(cart.Items) != 1 {
					t.Fatalf("expected 1 item, got %d", len(cart.Items))
				}
				if cart.Items[0].Quantity != 2 {
					t.Fatalf("expected quantity 2, got %d", cart.Items[0].Quantity)
				}
				return nil
			},
		}

		if _, err := NewService(runner, repo, nil).AddItem(context.Background(), uuid.New(), AddCartItemRequest{
			ProductID: productID, Quantity: 1,
		}); err != nil {
			t.Fatal("unexpected error:", err)
		}
	})

	t.Run("getOrCreateCart fails: returns error", func(t *testing.T) {
		repo := &MockCartRepository{
			GetByUserIDFn: func(_ context.Context, _ db.DBQE, _ uuid.UUID) (*model.Cart, error) {
				return nil, fmt.Errorf("db error")
			},
			CreateFn: func(_ context.Context, _ db.DBQE, _ uuid.UUID) (*model.Cart, error) {
				return nil, fmt.Errorf("create failed")
			},
		}

		if _, err := NewService(runner, repo, nil).AddItem(context.Background(), uuid.New(), AddCartItemRequest{
			ProductID: uuid.New(), Quantity: 1,
		}); err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("save fails: returns error", func(t *testing.T) {
		repo := &MockCartRepository{
			GetByUserIDFn: func(_ context.Context, _ db.DBQE, _ uuid.UUID) (*model.Cart, error) {
				return &model.Cart{Items: []model.CartItem{}}, nil
			},
			SaveFn: func(_ context.Context, _ db.DBQE, _ *model.Cart) error {
				return fmt.Errorf("save failed")
			},
		}

		if _, err := NewService(runner, repo, nil).AddItem(context.Background(), uuid.New(), AddCartItemRequest{
			ProductID: uuid.New(), Quantity: 1,
		}); err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}
