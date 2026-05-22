package order

import (
	"booky-backend/internal/model"
	"booky-backend/internal/trans"
	"booky-backend/pkg/database"
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
)

// ── Mocks ─────────────────────────────────────────────────────────────────────

type MockOrderRepository struct {
	CreateFn           func(ctx context.Context, db database.DBQE, order model.Order) (*model.Order, error)
	GetByIDFn          func(ctx context.Context, db database.DBQE, id uuid.UUID) (*model.Order, error)
	GetAllFn           func(ctx context.Context, db database.DBQE, q *trans.PaginationQuery) ([]*model.Order, *trans.Page, error)
	TransitionStatusFn func(ctx context.Context, db database.DBQE, id uuid.UUID, from, to model.OrderStatus) error
	UpdateTotalPriceFn func(ctx context.Context, db database.DBQE, orderID uuid.UUID, total int) error
}

func (m *MockOrderRepository) Create(ctx context.Context, db database.DBQE, order model.Order) (*model.Order, error) {
	if m.CreateFn == nil {
		panic("CreateFn is not set")
	}
	return m.CreateFn(ctx, db, order)
}
func (m *MockOrderRepository) UpdateTotalPrice(ctx context.Context, db database.DBQE, orderID uuid.UUID, total int) error {
	if m.UpdateTotalPriceFn == nil {
		panic("UpdateTotalPriceFn is not set")
	}
	return m.UpdateTotalPriceFn(ctx, db, orderID, total)
}
func (m *MockOrderRepository) GetByID(ctx context.Context, db database.DBQE, id uuid.UUID) (*model.Order, error) {
	if m.GetByIDFn == nil {
		panic("GetByIDFn is not set")
	}
	return m.GetByIDFn(ctx, db, id)
}
func (m *MockOrderRepository) GetAll(ctx context.Context, db database.DBQE, q *trans.PaginationQuery) ([]*model.Order, *trans.Page, error) {
	if m.GetAllFn == nil {
		panic("GetAllFn is not set")
	}
	return m.GetAllFn(ctx, db, q)
}
func (m *MockOrderRepository) TransitionStatus(ctx context.Context, db database.DBQE, id uuid.UUID, from, to model.OrderStatus) error {
	if m.TransitionStatusFn == nil {
		panic("TransitionStatusFn is not set")
	}
	return m.TransitionStatusFn(ctx, db, id, from, to)
}

type MockRunner struct {
	WithTxFn func(ctx context.Context, fn func(tx database.DBQE) error) error
	DBFn     func() database.DBQE
}

func (m *MockRunner) WithTx(ctx context.Context, fn func(tx database.DBQE) error) error {
	if m.WithTxFn == nil {
		panic("WithTxFn is not set")
	}
	return m.WithTxFn(ctx, fn)
}
func (m *MockRunner) DB() database.DBQE {
	if m.DBFn == nil {
		panic("DBFn is not set")
	}
	return m.DBFn()
}

// ── Helpers ───────────────────────────────────────────────────────────────────

func execTx(_ context.Context, fn func(database.DBQE) error) error { return fn(nil) }
func noDB() database.DBQE                                          { return nil }

// ── TestGetByID ───────────────────────────────────────────────────────────────

func TestGetByID(t *testing.T) {
	runner := &MockRunner{DBFn: noDB}

	t.Run("success: returns order", func(t *testing.T) {
		orderID := uuid.New()
		repo := &MockOrderRepository{
			GetByIDFn: func(_ context.Context, _ database.DBQE, id uuid.UUID) (*model.Order, error) {
				return &model.Order{ID: id}, nil
			},
		}

		order, err := NewService(runner, repo).GetByID(context.Background(), orderID)
		if err != nil {
			t.Fatal("unexpected error:", err)
		}
		if order.ID != orderID {
			t.Fatalf("expected ID %v, got %v", orderID, order.ID)
		}
	})

	t.Run("not found: returns error", func(t *testing.T) {
		repo := &MockOrderRepository{
			GetByIDFn: func(_ context.Context, _ database.DBQE, _ uuid.UUID) (*model.Order, error) {
				return nil, fmt.Errorf("not found")
			},
		}

		if _, err := NewService(runner, repo).GetByID(context.Background(), uuid.New()); err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

// ── TestGetAll ────────────────────────────────────────────────────────────────

func TestGetAll(t *testing.T) {
	runner := &MockRunner{DBFn: noDB}

	t.Run("success: returns orders and page", func(t *testing.T) {
		orders := []*model.Order{{ID: uuid.New()}, {ID: uuid.New()}}
		page := &trans.Page{Total: 2}
		repo := &MockOrderRepository{
			GetAllFn: func(_ context.Context, _ database.DBQE, _ *trans.PaginationQuery) ([]*model.Order, *trans.Page, error) {
				return orders, page, nil
			},
		}

		got, gotPage, err := NewService(runner, repo).GetAll(context.Background(), &trans.PaginationQuery{})
		if err != nil {
			t.Fatal("unexpected error:", err)
		}
		if len(got) != 2 || gotPage.Total != 2 {
			t.Fatalf("unexpected result: %+v %+v", got, gotPage)
		}
	})

	t.Run("repo error: returns error", func(t *testing.T) {
		repo := &MockOrderRepository{
			GetAllFn: func(_ context.Context, _ database.DBQE, _ *trans.PaginationQuery) ([]*model.Order, *trans.Page, error) {
				return nil, nil, fmt.Errorf("db error")
			},
		}

		if _, _, err := NewService(runner, repo).GetAll(context.Background(), &trans.PaginationQuery{}); err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

// ── TestCancel ────────────────────────────────────────────────────────────────

func TestCancel(t *testing.T) {
	// NOTE: Cancel passes s.tx.DB() to TransitionStatus instead of the tx
	// received from WithTx — this is likely a bug. Tests reflect current behavior.
	runner := &MockRunner{WithTxFn: execTx, DBFn: noDB}

	t.Run("success: transitions pending → cancelled", func(t *testing.T) {
		orderID := uuid.New()
		repo := &MockOrderRepository{
			TransitionStatusFn: func(_ context.Context, _ database.DBQE, id uuid.UUID, from, to model.OrderStatus) error {
				if id != orderID {
					t.Fatalf("expected order ID %v, got %v", orderID, id)
				}
				if from != model.OrderStatusPending || to != model.OrderStatusCancelled {
					t.Fatalf("unexpected transition: %v → %v", from, to)
				}
				return nil
			},
		}

		if err := NewService(runner, repo).Cancel(context.Background(), orderID); err != nil {
			t.Fatal("unexpected error:", err)
		}
	})

	t.Run("transition fails: returns error", func(t *testing.T) {
		repo := &MockOrderRepository{
			TransitionStatusFn: func(_ context.Context, _ database.DBQE, _ uuid.UUID, _, _ model.OrderStatus) error {
				return fmt.Errorf("invalid transition")
			},
		}

		if err := NewService(runner, repo).Cancel(context.Background(), uuid.New()); err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

// ── TestConfirm ───────────────────────────────────────────────────────────────

func TestConfirm(t *testing.T) {
	// NOTE: same bug as Cancel — uses s.tx.DB() instead of tx inside WithTx.
	runner := &MockRunner{WithTxFn: execTx, DBFn: noDB}

	t.Run("success: transitions pending → confirmed", func(t *testing.T) {
		orderID := uuid.New()
		repo := &MockOrderRepository{
			TransitionStatusFn: func(_ context.Context, _ database.DBQE, id uuid.UUID, from, to model.OrderStatus) error {
				if id != orderID {
					t.Fatalf("expected order ID %v, got %v", orderID, id)
				}
				if from != model.OrderStatusPending || to != model.OrderStatusConfirmed {
					t.Fatalf("unexpected transition: %v → %v", from, to)
				}
				return nil
			},
		}

		if err := NewService(runner, repo).Confirm(context.Background(), orderID); err != nil {
			t.Fatal("unexpected error:", err)
		}
	})

	t.Run("transition fails: returns error", func(t *testing.T) {
		repo := &MockOrderRepository{
			TransitionStatusFn: func(_ context.Context, _ database.DBQE, _ uuid.UUID, _, _ model.OrderStatus) error {
				return fmt.Errorf("invalid transition")
			},
		}

		if err := NewService(runner, repo).Confirm(context.Background(), uuid.New()); err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}
