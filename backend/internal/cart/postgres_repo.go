package cart

import (
	"booky-backend/internal/model"
	"booky-backend/pkg/database"
	"context"

	"github.com/google/uuid"
)

type PostgresRepository struct{}

func NewPostgresRepository() CartRepository {
	return &PostgresRepository{}
}

func (r *PostgresRepository) Create(ctx context.Context, qe database.QueryExecutor, userID uuid.UUID) (*model.Cart, error) {
	var cart model.Cart
	err := qe.QueryRow(ctx,
		`INSERT INTO carts (user_id, created_at, updated_at)
		 VALUES ($1, now(), now()) RETURNING id, created_at, updated_at`,
		userID,
	).Scan(&cart.ID, &cart.CreatedAt, &cart.UpdatedAt)

	if err != nil {
		return nil, database.MapError(err)
	}

	return &cart, nil
}

func (r *PostgresRepository) GetByUserID(ctx context.Context, qe database.QueryExecutor, userID uuid.UUID) (*model.Cart, error) {
	var cart model.Cart
	err := qe.QueryRow(ctx,
		`SELECT id, user_id, created_at, updated_at FROM carts WHERE user_id=$1`,
		userID,
	).Scan(&cart.ID, &cart.UserID, &cart.CreatedAt, &cart.UpdatedAt)

	if err != nil {
		return nil, database.MapError(err)
	}

	rows, err := qe.Query(ctx,
		`SELECT product_id, quantity FROM cart_items WHERE cart_id=$1`,
		cart.ID,
	)

	if err != nil {
		return nil, database.MapError(err)
	}

	for rows.Next() {
		var item model.CartItem
		err := rows.Scan(&item.ProductID, &item.Quantity)
		if err != nil {
			return nil, database.MapError(err)
		}
		cart.Items = append(cart.Items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, database.MapError(err)
	}

	return &cart, nil
}

func (r *PostgresRepository) Save(ctx context.Context, qe database.QueryExecutor, cart *model.Cart) error {
	// 1. Lock cart row (REAL race protection)
	_, err := qe.Exec(ctx, `
		SELECT id 
		FROM carts 
		WHERE id = $1 
		FOR UPDATE
	`, cart.ID)

	if err != nil {
		return database.MapError(err)
	}

	// 3. Delete old items (safe because locked)
	_, err = qe.Exec(ctx, `
		DELETE FROM cart_items 
		WHERE cart_id = $1
	`, cart.ID)

	if err != nil {
		return database.MapError(err)
	}

	// 4. Insert fresh snapshot
	for _, item := range cart.Items {
		_, err = qe.Exec(ctx, `
			INSERT INTO cart_items (cart_id, product_id, quantity)
			VALUES ($1, $2, $3)
		`,
			cart.ID,
			item.ProductID,
			item.Quantity,
		)

		if err != nil {
			return database.MapError(err)
		}
	}

	return nil
}

func (r *PostgresRepository) Empty(ctx context.Context, qe database.QueryExecutor, userID uuid.UUID) error {
	_, err := qe.Exec(ctx,
		`DELETE FROM cart_items WHERE cart_id IN (SELECT id FROM carts WHERE user_id=$1)`,
		userID,
	)
	if err != nil {
		return err
	}

	return nil
}
