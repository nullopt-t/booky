package order

import (
	"booky-backend/internal/db"
	"booky-backend/internal/domain"
	"booky-backend/internal/utils"
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
)

var (
	ErrInDatabase          = errors.New("database error")
	ErrNoItems             = errors.New("no items in order")
	ErrInsufficientQuanity = errors.New("insufficient quantity")
	ErrInvalidQuantity     = errors.New("invalid quantity")
	ErrProductNotFound     = errors.New("product not found")
	ErrInvalidProductID    = errors.New("invalid product ID")
	ErrOrderNotFound       = errors.New("order not found")
	ErrOrderNotPending     = errors.New("order is not pending")
)

type PostgresRepo struct {
	db *db.DB
}

func NewPostgresRepo(db *db.DB) *PostgresRepo {
	return &PostgresRepo{
		db,
	}
}

func (r *PostgresRepo) Create(ctx context.Context, order CreateOrderRequest) (*domain.Order, error) {
	if len(order.Items) == 0 {
		return nil, ErrNoItems
	}

	tx, err := r.db.GetPool().Begin(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = tx.Rollback(ctx) // no-op if already committed
	}()

	var orderID string
	var createdOrder domain.Order
	err = tx.QueryRow(ctx, `INSERT INTO orders(total_price) VALUES ($1) RETURNING id`, 0).Scan(&orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to create new order: %w", err)
	}

	var totalPrice int
	for _, item := range order.Items {
		var stock, price int
		err := tx.QueryRow(ctx, "SELECT stock, price FROM products WHERE id = $1 FOR UPDATE", item.ProductID).Scan(&stock, &price)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, ErrProductNotFound
			}
			return nil, fmt.Errorf("failed to get product stock and price: %w", err)
		}
		if stock < item.Quantity {
			return nil, ErrInsufficientQuanity
		}

		_, err = tx.Exec(ctx, "UPDATE products SET stock = stock - $1 WHERE id = $2", item.Quantity, item.ProductID)
		if err != nil {
			return nil, fmt.Errorf("failed to update product stock: %w", err)
		}

		_, err = tx.Exec(ctx, `INSERT INTO order_items (order_id, product_id, quantity, purchase_price) VALUES ($1, $2, $3, $4)`, orderID, item.ProductID, item.Quantity, price)
		if err != nil {
			return nil, fmt.Errorf("failed to insert order item: %w", err)
		}

		createdOrder.Items = append(createdOrder.Items, domain.OrderItem{
			ProductID:     item.ProductID,
			Quantity:      item.Quantity,
			PurchasePrice: price,
		})

		totalPrice += price * item.Quantity
	}

	err = tx.QueryRow(ctx, "UPDATE orders SET total_price = $1 WHERE id = $2 RETURNING id, status, total_price, created_at, updated_at", totalPrice, orderID).Scan(&createdOrder.ID, &createdOrder.Status, &createdOrder.TotalPrice, &createdOrder.CreatedAt, &createdOrder.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to update order total price: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &createdOrder, nil
}

func (r *PostgresRepo) Cancel(ctx context.Context, orderID string) error {
	tx, err := r.db.GetPool().Begin(ctx)
	if err != nil {
		log.Println(err)
		return ErrInDatabase
	}

	defer tx.Rollback(ctx)

	var status domain.OrderStatus
	err = tx.QueryRow(ctx, "SELECT status FROM orders WHERE id = $1 FOR UPDATE", orderID).Scan(&status)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrOrderNotFound
		}
		log.Println(err)
		return ErrInDatabase
	}

	if status != domain.OrderStatusPending {
		return ErrOrderNotPending
	}

	_, err = tx.Exec(ctx, "UPDATE products SET stock = stock + (SELECT quantity FROM order_items WHERE order_id = $1 AND product_id = products.id) WHERE id IN (SELECT product_id FROM order_items WHERE order_id = $1)", orderID)
	if err != nil {
		return fmt.Errorf("failed to revert order items: %w", err)
	}

	_, err = tx.Exec(ctx, "UPDATE orders SET status = $1 WHERE id = $2", domain.OrderStatusCancelled, orderID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrOrderNotFound
		}
		return fmt.Errorf("failed to cancel order: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *PostgresRepo) Confirm(ctx context.Context, orderID string) error {
	tx, err := r.db.GetPool().Begin(ctx)
	if err != nil {
		log.Println(err)
		return ErrInDatabase
	}

	defer tx.Rollback(ctx)

	var status domain.OrderStatus
	err = tx.QueryRow(ctx, "SELECT status FROM orders WHERE id = $1 FOR UPDATE", orderID).Scan(&status)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrOrderNotFound
		}
		log.Println(err)
		return ErrInDatabase
	}

	if status != domain.OrderStatusPending {
		return ErrOrderNotPending
	}

	_, err = tx.Exec(ctx, "UPDATE orders SET status = $1 WHERE id = $2", domain.OrderStatusConfirmed, orderID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrOrderNotFound
		}
		return fmt.Errorf("failed to confirm order: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *PostgresRepo) GetByID(ctx context.Context, id string) (*domain.Order, error) {
	var order domain.Order
	err := r.db.GetPool().QueryRow(ctx, `
    SELECT
      o.id,
      o.status,
      o.total_price,
      o.created_at,
	  o.updated_at
    FROM orders o
    WHERE o.id = $1
    `, id).Scan(&order.ID, &order.Status, &order.TotalPrice, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrOrderNotFound
		}
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	var items = make([]domain.OrderItem, 0)
	rows, err := r.db.GetPool().Query(ctx, "SELECT order_id, quantity, purchase_price FROM order_items WHERE order_id = $1 ORDER BY created_at DESC", id)
	if err != nil {
		return nil, fmt.Errorf("failed to get order items: %w", err)
	}

	for rows.Next() {
		var item domain.OrderItem
		if err := rows.Scan(&item.ProductID, &item.Quantity, &item.PurchasePrice); err != nil {
			return nil, fmt.Errorf("failed to scan order item: %w", err)
		}
		items = append(items, item)
	}
	defer rows.Close()
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate over order items: %w", err)
	}

	order.Items = items
	return &order, nil
}

func (r *PostgresRepo) GetAll(ctx context.Context, q utils.PaginationQuery) ([]*OrderResponse, error) {
	// build base query with pagination
	sql := `
    SELECT
      o.id,
      o.status,
      o.total_price,
      o.created_at,
      oi.id AS item_id,
      oi.product_id,
      oi.quantity,
      oi.purchase_price
    FROM orders o
    JOIN order_items oi ON oi.order_id = o.id
    JOIN products p ON p.id = oi.product_id
    ORDER BY o.created_at DESC
    `
	// append LIMIT/OFFSET from q (assumes q.Limit and q.Offset int, 0 means no limit)
	args := []interface{}{}
	offset := (q.Page - 1) * q.Limit
	if q.Limit > 0 {
		sql += " LIMIT $1"
		args = append(args, q.Limit)
		if offset > 0 {
			sql += " OFFSET $2"
			args = append(args, offset)
		}
	} else if offset > 0 { // offset without limit is uncommon but handled
		sql += " OFFSET $1"
		args = append(args, offset)
	}

	rows, err := r.db.GetPool().Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ordersMap := make(map[string]*OrderResponse)
	for rows.Next() {
		var (
			orderID       string
			status        domain.OrderStatus
			totalPrice    int
			createdAt     time.Time
			itemID        string
			productID     string
			quantity      int
			purchasePrice int
		)

		if err := rows.Scan(
			&orderID,
			&status,
			&totalPrice,
			&createdAt,
			&itemID,
			&productID,
			&quantity,
			&purchasePrice,
		); err != nil {
			return nil, err
		}

		o, ok := ordersMap[orderID]
		if !ok {
			o = &OrderResponse{
				ID:         orderID,
				Status:     status,
				TotalPrice: totalPrice,
				CreatedAt:  createdAt,
				Items:      []OrderItemResponse{},
			}
			ordersMap[orderID] = o
		}

		o.Items = append(o.Items, OrderItemResponse{
			ProductID:     productID,
			Quantity:      quantity,
			PurchasePrice: purchasePrice,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// convert map to slice preserving order of created_at DESC
	results := make([]*OrderResponse, 0, len(ordersMap))
	// if you need stable ordering, query order IDs first then fetch items separately;
	// here we iterate the map (ordering not guaranteed). For small result sets this is fine.
	for _, o := range ordersMap {
		results = append(results, o)
	}

	return results, nil
}
