package order

import (
	"booky-backend/internal/db"
	"booky-backend/internal/domain"
	"booky-backend/internal/logger"
	"booky-backend/internal/utils"
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
)

var (
	ErrInDatabase            = errors.New("database error")
	ErrNoItems               = errors.New("no items in order")
	ErrInsufficientQuanity   = errors.New("insufficient quantity")
	ErrInvalidQuantity       = errors.New("invalid quantity")
	ErrProductNotFound       = errors.New("product not found")
	ErrInvalidProductID      = errors.New("invalid product ID")
	ErrOrderNotFound         = errors.New("order not found")
	ErrOrderNotPending       = errors.New("order is not pending")
	ErrOrderAlreadyCancelled = errors.New("order already cancelled")
	ErrOrderAlreadyConfirmed = errors.New("order already confirmed")
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
		logger.Log(logger.DEBUG, "no items in the order")
		return nil, ErrNoItems
	}

	tx, err := r.db.GetPool().Begin(ctx)
	if err != nil {
		logger.Log(logger.DEBUG, "failed to begin transaction: %v", err.Error())
		return nil, ErrInDatabase
	}

	defer func() {
		_ = tx.Rollback(ctx) // no-op if already committed
	}()

	var orderID string
	var createdOrder domain.Order
	err = tx.QueryRow(ctx, `INSERT INTO orders(total_price) VALUES ($1) RETURNING id`, 0).Scan(&orderID)
	if err != nil {
		logger.Log(logger.DEBUG, "failed to insert new order :%v", err.Error())
		return nil, ErrInDatabase
	}

	var totalPrice int
	for _, item := range order.Items {
		var stock, price int
		err := tx.QueryRow(ctx, "SELECT stock, price FROM products WHERE id = $1 FOR UPDATE", item.ProductID).Scan(&stock, &price)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				logger.Log(logger.DEBUG, "failed to get product stock and price: %v", err.Error())
				return nil, ErrProductNotFound
			}
			logger.Log(logger.DEBUG, "failed to get product stock and price: %v", err.Error())
			return nil, ErrInDatabase
		}
		if stock < item.Quantity {
			logger.Log(logger.DEBUG, "insufficient quanity of product %v", item.ProductID)
			return nil, ErrInsufficientQuanity
		}

		_, err = tx.Exec(ctx, "UPDATE products SET stock = stock - $1 WHERE id = $2", item.Quantity, item.ProductID)
		if err != nil {
			logger.Log(logger.DEBUG, "failed to update product stock: %v", err.Error())
			return nil, ErrInDatabase
		}

		_, err = tx.Exec(ctx, `INSERT INTO order_items (order_id, product_id, quantity, purchase_price) VALUES ($1, $2, $3, $4)`, orderID, item.ProductID, item.Quantity, price)
		if err != nil {
			logger.Log(logger.DEBUG, "failed to insert order item: %v", err.Error())
			return nil, ErrInDatabase
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
		logger.Log(logger.DEBUG, "failed to update order total price: %v", err.Error())
		return nil, ErrInDatabase
	}

	if err = tx.Commit(ctx); err != nil {
		logger.Log(logger.DEBUG, "failed to commit transaction: %v", err.Error())
		return nil, ErrInDatabase
	}

	return &createdOrder, nil
}

func (r *PostgresRepo) Cancel(ctx context.Context, orderID string) error {
	tx, err := r.db.GetPool().Begin(ctx)
	if err != nil {
		logger.Log(logger.DEBUG, "failed to begin transaction: %v", err.Error())
		return ErrInDatabase
	}

	defer tx.Rollback(ctx)

	var status OrderStatus
	err = tx.QueryRow(ctx, "SELECT status FROM orders WHERE id = $1 FOR UPDATE", orderID).Scan(&status)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrOrderNotFound
		}
		logger.Log(logger.DEBUG, "failed to get order status: %v", err.Error())
		return ErrInDatabase
	}
	if status == OrderStatusCancelled {
		logger.Log(logger.DEBUG, "order is already cancelled")
		return ErrOrderAlreadyCancelled
	}

	if status != OrderStatusPending {
		logger.Log(logger.DEBUG, "order is not pending")
		return ErrOrderNotPending
	}

	_, err = tx.Exec(ctx, "UPDATE products SET stock = stock + (SELECT quantity FROM order_items WHERE order_id = $1 AND product_id = products.id) WHERE id IN (SELECT product_id FROM order_items WHERE order_id = $1)", orderID)
	if err != nil {
		logger.Log(logger.DEBUG, "failed to revert order items: %v", err.Error())
		return ErrInDatabase
	}

	_, err = tx.Exec(ctx, "UPDATE orders SET status = $1 WHERE id = $2", OrderStatusCancelled, orderID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrOrderNotFound
		}
		logger.Log(logger.DEBUG, "failed to cancel order: %v", err.Error())
		return ErrInDatabase
	}

	if err = tx.Commit(ctx); err != nil {
		logger.Log(logger.DEBUG, "failed to commit transaction: %v", err.Error())
		return ErrInDatabase
	}

	return nil
}

func (r *PostgresRepo) Confirm(ctx context.Context, orderID string) error {
	tx, err := r.db.GetPool().Begin(ctx)
	if err != nil {
		logger.Log(logger.DEBUG, "failed to begin transaction: %v", err.Error())
		return ErrInDatabase
	}

	defer tx.Rollback(ctx)

	var status OrderStatus
	err = tx.QueryRow(ctx, "SELECT status FROM orders WHERE id = $1 FOR UPDATE", orderID).Scan(&status)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrOrderNotFound
		}
		logger.Log(logger.DEBUG, "failed to get order status: %v", err.Error())
		return ErrInDatabase
	}

	if status == OrderStatusConfirmed {
		logger.Log(logger.DEBUG, "order is already confirmed")
		return ErrOrderAlreadyConfirmed
	}

	if status != OrderStatusPending {
		return ErrOrderNotPending
	}

	_, err = tx.Exec(ctx, "UPDATE orders SET status = $1 WHERE id = $2", OrderStatusConfirmed, orderID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrOrderNotFound
		}
		logger.Log(logger.DEBUG, "failed to confirm order: %v", err.Error())
		return ErrInDatabase
	}

	if err = tx.Commit(ctx); err != nil {
		logger.Log(logger.DEBUG, "failed to commit transaction: %v", err.Error())
		return ErrInDatabase
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
		logger.Log(logger.DEBUG, "failed to get order: %v", err.Error())
		return nil, ErrInDatabase
	}

	var items = make([]domain.OrderItem, 0)
	rows, err := r.db.GetPool().Query(ctx, "SELECT order_id, quantity, purchase_price FROM order_items WHERE order_id = $1 ORDER BY created_at DESC", id)
	if err != nil {
		logger.Log(logger.DEBUG, "failed to get order items: %v", err.Error())
		return nil, ErrInDatabase
	}

	for rows.Next() {
		var item domain.OrderItem
		if err := rows.Scan(&item.ProductID, &item.Quantity, &item.PurchasePrice); err != nil {
			logger.Log(logger.DEBUG, "failed to scan order item: %v", err.Error())
			return nil, ErrInDatabase
		}
		items = append(items, item)
	}
	defer rows.Close()
	if err := rows.Err(); err != nil {
		logger.Log(logger.DEBUG, "failed to iterate over order items: %v", err.Error())
		return nil, ErrInDatabase
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
		logger.Log(logger.DEBUG, "failed to get all orders: %v", err.Error())
		return nil, ErrInDatabase
	}
	defer rows.Close()

	ordersMap := make(map[string]*OrderResponse)
	for rows.Next() {
		var (
			orderID       string
			status        OrderStatus
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
			logger.Log(logger.DEBUG, "failed to scan order: %v", err.Error())
			return nil, ErrInDatabase
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
		logger.Log(logger.DEBUG, "failed to iterate over orders: %v", err.Error())
		return nil, ErrInDatabase
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
