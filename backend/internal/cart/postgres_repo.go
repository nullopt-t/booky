package cart

import (
	"booky-backend/internal/db"
	"booky-backend/internal/model"
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type PostgresRepository struct{}

func NewPostgresRepository() CartRepository {
	return &PostgresRepository{}
}

func (r *PostgresRepository) Create(ctx context.Context, qe db.DBQE, userID uuid.UUID) (*model.Cart, error) {
	var cart model.Cart
	err := qe.QueryRow(ctx,
		`INSERT INTO carts (user_id, created_at, updated_at)
		 VALUES ($1, now(), now()) RETURNING id, created_at, updated_at`,
		userID,
	).Scan(&cart.ID, &cart.CreatedAt, &cart.UpdatedAt)

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == db.UniqueViolation {
				return nil, fmt.Errorf("%w : %v", ErrCartAlreadyExist, err)
			}
		}
		return nil, fmt.Errorf("%w : %v", ErrDatabaseFailure, err)
	}

	return &cart, nil
}

func (r *PostgresRepository) GetByUserID(ctx context.Context, qe db.DBQE, userID uuid.UUID) (*model.Cart, error) {
	var cart model.Cart
	err := qe.QueryRow(ctx,
		`SELECT id, user_id, created_at, updated_at FROM carts WHERE user_id=$1`,
		userID,
	).Scan(&cart.ID, &cart.UserID, &cart.CreatedAt, &cart.UpdatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("%w : %v", ErrCartNotFound, err)
		}
		return nil, fmt.Errorf("%w : %v", ErrDatabaseFailure, err)
	}

	rows, err := qe.Query(ctx,
		`SELECT product_id, quantity FROM cart_items WHERE cart_id=$1`,
		cart.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("%w : %v", ErrDatabaseFailure, err)
	}

	for rows.Next() {
		var item model.CartItem
		err := rows.Scan(&item.ProductID, &item.Quantity)
		if err != nil {
			return nil, fmt.Errorf("%w : %v", ErrDatabaseFailure, err)
		}
		cart.Items = append(cart.Items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w : %v", ErrDatabaseFailure, err)
	}

	return &cart, nil
}

func (r *PostgresRepository) Save(ctx context.Context, qe db.DBQE, cart *model.Cart) error {
	// 1. Lock cart row (REAL race protection)
	_, err := qe.Exec(ctx, `
		SELECT id 
		FROM carts 
		WHERE id = $1 
		FOR UPDATE
	`, cart.ID)

	if err != nil {
		return fmt.Errorf("%w : %v", ErrDatabaseFailure, err)
	}

	// 2. Update cart timestamp
	_, err = qe.Exec(ctx, `
		UPDATE carts 
		SET updated_at = now()
		WHERE id = $1
	`, cart.ID)

	if err != nil {
		return fmt.Errorf("%w : %v", ErrDatabaseFailure, err)
	}

	// 3. Delete old items (safe because locked)
	_, err = qe.Exec(ctx, `
		DELETE FROM cart_items 
		WHERE cart_id = $1
	`, cart.ID)

	if err != nil {
		return fmt.Errorf("%w : %v", ErrDatabaseFailure, err)
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
			return fmt.Errorf("%w : %v", ErrDatabaseFailure, err)
		}
	}

	return nil
}

func (r *PostgresRepository) Empty(ctx context.Context, qe db.DBQE, userID uuid.UUID) error {
	_, err := qe.Exec(ctx,
		`DELETE FROM cart_items WHERE cart_id IN (SELECT id FROM carts WHERE user_id=$1)`,
		userID,
	)
	if err != nil {
		return fmt.Errorf("%w : %v", ErrDatabaseFailure, err)
	}
	return nil
}
