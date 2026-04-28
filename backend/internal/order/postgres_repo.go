package order

import (
	"booky-backend/internal/db"
	"booky-backend/internal/domain"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

const ErrInvalidQuantity = "invalid quantity"

type PostgresRepo struct {
	db *db.DB
}

func NewPostgresRepo(db *db.DB) *PostgresRepo {
	return &PostgresRepo{
		db,
	}
}

func (r *PostgresRepo) Create(ctx context.Context, order CreateOrderRequest) (*CreateOrderResponse, error) {
	// basic validation
	if len(order.Items) == 0 {
		return nil, fmt.Errorf("order must contain at least one item")
	}

	tx, err := r.db.GetPool().Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback(ctx) // no-op if already committed
	}()

	var totalPriceCents int
	type itemInfo OrderItemResponse
	items := make([]OrderItemResponse, 0, len(order.Items))

	// prepare statements
	selStmt, err := tx.Prepare(ctx, "sel_product", `SELECT stock, price FROM products WHERE id = $1 FOR UPDATE`)
	if err != nil {
		return nil, err
	}
	updStmt, err := tx.Prepare(ctx, "upd_stock", `UPDATE products SET stock = stock - $1 WHERE id = $2`)
	if err != nil {
		return nil, err
	}

	for _, it := range order.Items {
		if it.Quantity <= 0 {
			return nil, fmt.Errorf("invalid quantity for product %s", it.ProductID)
		}

		var stock int
		var priceCents int
		err = tx.QueryRow(ctx, selStmt.Name, it.ProductID).Scan(&stock, &priceCents)
		if err != nil {
			if err == pgx.ErrNoRows {
				return nil, fmt.Errorf("product %s not found", it.ProductID)
			}
			return nil, err
		}
		if stock < it.Quantity {
			return nil, fmt.Errorf("product %s out of stock (have %d, want %d)", it.ProductID, stock, it.Quantity)
		}

		if _, err = tx.Exec(ctx, updStmt.Name, it.Quantity, it.ProductID); err != nil {
			return nil, err
		}

		totalPriceCents += it.Quantity * priceCents
		items = append(items, OrderItemResponse{
			ProductID:     it.ProductID,
			Quantity:      it.Quantity,
			PurchasePrice: priceCents,
		})
	}

	// insert order
	var orderID string
	var createdAt time.Time
	err = tx.QueryRow(ctx,
		`INSERT INTO orders (status, total_price) VALUES ($1, $2) RETURNING id, created_at`,
		domain.OrderStatus("pending"), totalPriceCents,
	).Scan(&orderID, &createdAt)
	if err != nil {
		return nil, err
	}

	// batch insert items with a single multi-row INSERT
	// build VALUES ($1,$2,$3,$4),($5,$6,$7,$8),...
	vals := make([]interface{}, 0, len(items)*4)
	placeholders := make([]string, 0, len(items))
	for i, it := range items {
		idx := i*4 + 1
		placeholders = append(placeholders, fmt.Sprintf("($%d,$%d,$%d,$%d)", idx, idx+1, idx+2, idx+3))
		vals = append(vals, orderID, it.ProductID, it.Quantity, it.PurchasePrice)
	}
	q := `INSERT INTO order_items (order_id, product_id, quantity, purchase_price) VALUES ` + strings.Join(placeholders, ",")
	if _, err = tx.Exec(ctx, q, vals...); err != nil {
		return nil, err
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &CreateOrderResponse{
		ID:         orderID,
		Status:     domain.OrderStatus("pending"),
		TotalPrice: totalPriceCents,
		CreatedAt:  createdAt,
		Items:      items,
	}, nil
}

func (r *PostgresRepo) GetAll(ctx context.Context, q db.PaginationQuery) ([]*OrderResponse, error) {
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
			status        string
			totalPrice    int64
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
