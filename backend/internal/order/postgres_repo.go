package order

import (
	"booky-backend/internal/model"
	"booky-backend/pkg/api"
	"booky-backend/pkg/database"
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type PostgresRepo struct {
}

func NewPostgresRepository() OrderRepository {
	return &PostgresRepo{}
}

func (r *PostgresRepo) Create(ctx context.Context, db database.QueryExecutor, order model.Order) (*model.Order, error) {
	// create order
	var createdOrder model.Order
	err := db.QueryRow(ctx, `INSERT INTO orders(total_price) VALUES ($1) RETURNING id, status, total_price, created_at, updated_at`,
		order.TotalPrice).Scan(&createdOrder.ID, &createdOrder.Status, &createdOrder.TotalPrice, &createdOrder.CreatedAt, &createdOrder.UpdatedAt)
	if err != nil {
		return nil, ErrInDatabase
	}

	// insert the order items
	for _, item := range order.Items {
		_, err = db.Exec(ctx, `INSERT INTO order_items (order_id, product_id, quantity, purchase_price) VALUES ($1, $2, $3, $4)`,
			createdOrder.ID, item.ProductID, item.Quantity, item.PurchasePrice)
		if err != nil {
			return nil, ErrInDatabase
		}

		createdOrder.Items = append(createdOrder.Items, model.OrderItem{
			ProductID:     item.ProductID,
			Quantity:      item.Quantity,
			PurchasePrice: item.PurchasePrice,
		})
	}

	return &createdOrder, nil
}

func (r *PostgresRepo) GetByID(ctx context.Context, db database.QueryExecutor, orderID uuid.UUID) (*model.Order, error) {
	var order model.Order
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
		return nil, ErrInDatabase
	}

	var items = make([]model.OrderItem, 0)
	rows, err := db.Query(ctx, "SELECT order_id, quantity, purchase_price FROM order_items WHERE order_id = $1 ORDER BY created_at DESC", orderID)
	if err != nil {
		return nil, ErrInDatabase
	}

	for rows.Next() {
		var item model.OrderItem
		if err := rows.Scan(&item.ProductID, &item.Quantity, &item.PurchasePrice); err != nil {
			return nil, ErrInDatabase
		}
		items = append(items, item)
	}
	defer rows.Close()
	if err := rows.Err(); err != nil {
		return nil, ErrInDatabase
	}

	order.Items = items
	return &order, nil
}

func (r *PostgresRepo) GetAll(ctx context.Context, db database.QueryExecutor, q *api.PageQuery) ([]*model.Order, *api.Page, error) {
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
		return nil, nil, ErrInDatabase
	}
	defer rows.Close()

	ordersMap := make(map[uuid.UUID]*model.Order)
	for rows.Next() {
		var order model.Order
		var item model.OrderItem

		if err := rows.Scan(
			&order.ID,
			&order.Status,
			&order.TotalPrice,
			&order.CreatedAt,
			&item.ProductID,
			&item.Quantity,
			&item.PurchasePrice,
		); err != nil {
			return nil, nil, ErrInDatabase
		}

		o, ok := ordersMap[order.ID]
		if !ok {
			o = &model.Order{
				ID:         order.ID,
				Status:     order.Status,
				TotalPrice: order.TotalPrice,
				CreatedAt:  order.CreatedAt,
				Items:      []model.OrderItem{},
			}
			ordersMap[order.ID] = o
		}

		o.Items = append(o.Items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, nil, ErrInDatabase
	}

	// convert map to slice preserving order of created_at DESC
	results := make([]*model.Order, 0, len(ordersMap))
	for _, o := range ordersMap {
		results = append(results, o)
	}

	// query the orders count
	var count int
	err = db.QueryRow(ctx, "SELECT COUNT(*) FROM orders").Scan(&count)
	if err != nil {
		return nil, nil, ErrInDatabase
	}

	resultPage := &api.Page{
		Index: q.Page,
		Limit: q.Limit,
		Total: count,
	}

	return results, resultPage, nil
}

func (r *PostgresRepo) TransitionStatus(
	ctx context.Context,
	db database.QueryExecutor,
	orderID uuid.UUID,
	from,
	to model.OrderStatus,
) error {

	var currentStatus model.OrderStatus

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

func (r *PostgresRepo) UpdateTotalPrice(ctx context.Context, db database.QueryExecutor, orderID uuid.UUID, total int) error {
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
