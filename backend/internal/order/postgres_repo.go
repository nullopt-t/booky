package order

import (
	"booky-backend/internal/db"
	"booky-backend/internal/shared"
	"booky-backend/internal/trans"
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type PostgresRepo struct {
}

func NewPostgresRepo() OrderRepository {
	return &PostgresRepo{}
}

func (r *PostgresRepo) Create(ctx context.Context, db db.DBQE, order *CreateOrderRequest) (*Order, error) {
	var orderID string
	var createdOrder Order
	err := db.QueryRow(ctx, `INSERT INTO orders(total_price) VALUES ($1) RETURNING id`, 0).Scan(&orderID)
	if err != nil {
		shared.Log(shared.DEBUG, "failed to insert new order :%v", err.Error())
		return nil, ErrInDatabase
	}

	var totalPrice int
	for _, item := range order.Items {
		var stock, price int
		err := db.QueryRow(ctx, "SELECT stock, price FROM products WHERE id = $1 FOR UPDATE", item.ProductID).Scan(&stock, &price)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				shared.Log(shared.DEBUG, "failed to get product stock and price: %v", err.Error())
				return nil, ErrProductNotFound
			}
			shared.Log(shared.DEBUG, "failed to get product stock and price: %v", err.Error())
			return nil, ErrInDatabase
		}
		if stock < item.Quantity {
			shared.Log(shared.DEBUG, "insufficient quanity of product %v", item.ProductID)
			return nil, ErrInsufficientQuanity
		}

		_, err = db.Exec(ctx, "UPDATE products SET stock = stock - $1 WHERE id = $2", item.Quantity, item.ProductID)
		if err != nil {
			shared.Log(shared.DEBUG, "failed to update product stock: %v", err.Error())
			return nil, ErrInDatabase
		}

		_, err = db.Exec(ctx, `INSERT INTO order_items (order_id, product_id, quantity, purchase_price) VALUES ($1, $2, $3, $4)`, orderID, item.ProductID, item.Quantity, price)
		if err != nil {
			shared.Log(shared.DEBUG, "failed to insert order item: %v", err.Error())
			return nil, ErrInDatabase
		}

		createdOrder.Items = append(createdOrder.Items, OrderItem{
			ProductID:     item.ProductID,
			Quantity:      item.Quantity,
			PurchasePrice: price,
		})

		totalPrice += price * item.Quantity
	}

	err = db.QueryRow(ctx, "UPDATE orders SET total_price = $1 WHERE id = $2 RETURNING id, status, total_price, created_at, updated_at", totalPrice, orderID).Scan(&createdOrder.ID, &createdOrder.Status, &createdOrder.TotalPrice, &createdOrder.CreatedAt, &createdOrder.UpdatedAt)
	if err != nil {
		shared.Log(shared.DEBUG, "failed to update order total price: %v", err.Error())
		return nil, ErrInDatabase
	}

	return &createdOrder, nil
}

func (r *PostgresRepo) GetByID(ctx context.Context, db db.DBQE, orderID uuid.UUID) (*Order, error) {
	var order Order
	err := db.QueryRow(ctx, `
    SELECT
      o.id,
      o.status,
      o.total_price,
      o.created_at,
	  o.updated_at
    FROM orders o
    WHERE o.id = $1
    `, orderID).Scan(&order.ID, &order.Status, &order.TotalPrice, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrOrderNotFound
		}
		shared.Log(shared.DEBUG, "failed to get order: %v", err.Error())
		return nil, ErrInDatabase
	}

	var items = make([]OrderItem, 0)
	rows, err := db.Query(ctx, "SELECT order_id, quantity, purchase_price FROM order_items WHERE order_id = $1 ORDER BY created_at DESC", orderID)
	if err != nil {
		shared.Log(shared.DEBUG, "failed to get order items: %v", err.Error())
		return nil, ErrInDatabase
	}

	for rows.Next() {
		var item OrderItem
		if err := rows.Scan(&item.ProductID, &item.Quantity, &item.PurchasePrice); err != nil {
			shared.Log(shared.DEBUG, "failed to scan order item: %v", err.Error())
			return nil, ErrInDatabase
		}
		items = append(items, item)
	}
	defer rows.Close()
	if err := rows.Err(); err != nil {
		shared.Log(shared.DEBUG, "failed to iterate over order items: %v", err.Error())
		return nil, ErrInDatabase
	}

	order.Items = items
	return &order, nil
}

func (r *PostgresRepo) GetAll(ctx context.Context, db db.DBQE, q *trans.PaginationQuery) ([]*Order, *trans.Page, error) {
	// build base query with pagination
	sql := `
    SELECT
      o.id,
      o.status,
      o.total_price,
      o.created_at,
      oi.product_id,
      oi.quantity,
      oi.purchase_price
    FROM orders o
    JOIN order_items oi ON oi.order_id = o.id
    JOIN products p ON p.id = oi.product_id
	ORDER BY o.created_at DESC
    LIMIT $1 OFFSET $2
    `

	rows, err := db.Query(ctx, sql, q.Limit, (q.Page-1)*q.Limit)
	if err != nil {
		shared.Log(shared.DEBUG, "failed to get all orders: %v", err.Error())
		return nil, nil, ErrInDatabase
	}
	defer rows.Close()

	ordersMap := make(map[string]*Order)
	for rows.Next() {
		var order Order
		var item OrderItem

		if err := rows.Scan(
			&order.ID,
			&order.Status,
			&order.TotalPrice,
			&order.CreatedAt,
			&item.ProductID,
			&item.Quantity,
			&item.PurchasePrice,
		); err != nil {
			shared.Log(shared.DEBUG, "failed to scan order: %v", err.Error())
			return nil, nil, ErrInDatabase
		}

		o, ok := ordersMap[order.ID]
		if !ok {
			o = &Order{
				ID:         order.ID,
				Status:     order.Status,
				TotalPrice: order.TotalPrice,
				CreatedAt:  order.CreatedAt,
				Items:      []OrderItem{},
			}
			ordersMap[order.ID] = o
		}

		o.Items = append(o.Items, item)
	}

	if err := rows.Err(); err != nil {
		shared.Log(shared.DEBUG, "failed to iterate over orders: %v", err.Error())
		return nil, nil, ErrInDatabase
	}

	// convert map to slice preserving order of created_at DESC
	results := make([]*Order, 0, len(ordersMap))
	for _, o := range ordersMap {
		results = append(results, o)
	}

	// query the orders count
	var count int
	err = db.QueryRow(ctx, "SELECT COUNT(*) FROM orders").Scan(&count)
	if err != nil {
		return nil, nil, ErrInDatabase
	}

	resultPage := &trans.Page{
		Index: q.Page,
		Limit: q.Limit,
		Total: count,
	}

	return results, resultPage, nil
}

func (r *PostgresRepo) TransitionStatus(
	ctx context.Context,
	db db.DBQE,
	orderID uuid.UUID,
	from,
	to OrderStatus,
) error {

	var currentStatus OrderStatus

	err := db.QueryRow(ctx, `
        SELECT status
        FROM orders
        WHERE id = $1
    `, orderID).Scan(&currentStatus)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrOrderNotFound
		}

		return ErrInDatabase
	}

	// already in target state = idempotent success
	if currentStatus == to {
		return nil
	}

	// invalid transition
	if currentStatus != from {
		return ErrInvalidOrderTransition
	}

	_, err = db.Exec(ctx, `
        UPDATE orders
        SET status = $1
        WHERE id = $2
    `, to, orderID)

	if err != nil {
		return ErrInDatabase
	}

	return nil
}

func (r *PostgresRepo) UpdateTotalPrice(ctx context.Context, db db.DBQE, orderID uuid.UUID, total int) error {
	_, err := db.Exec(ctx, `
		UPDATE orders
		SET total_price = $1
		WHERE id = $2
	`, total, orderID)
	if err != nil {
		return ErrInDatabase
	}

	return nil
}
