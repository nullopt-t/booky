package cart

import (
	"booky-backend/internal/shared"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) CartRepository {
	return &PostgresRepository{
		db: db,
	}
}

func (r *PostgresRepository) Create(ctx context.Context, userID string) (*Cart, error) {
	var cart Cart
	err := r.db.QueryRow(ctx,
		`INSERT INTO carts (user_id, created_at, updated_at)
		 VALUES ($1, now(), now()) RETURNING id, user_id, created_at, updated_at`,
		userID,
	).Scan(&cart.ID, &cart.UserID, &cart.CreatedAt, &cart.UpdatedAt)

	if err != nil {
		shared.Log(shared.ERROR, err.Error())
		return nil, ErrInDatabase
	}

	return &cart, nil
}

func (r *PostgresRepository) GetByUserID(ctx context.Context, userID string) (*Cart, error) {
	var cart Cart
	err := r.db.QueryRow(ctx,
		`SELECT id, user_id, created_at, updated_at FROM carts WHERE user_id=$1`,
		userID,
	).Scan(&cart.ID, &cart.UserID, &cart.CreatedAt, &cart.UpdatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrCartNotFound
		}
		shared.Log(shared.ERROR, err.Error())
		return nil, ErrInDatabase
	}

	rows, err := r.db.Query(ctx,
		`SELECT item_id, quantity FROM cart_items WHERE cart_id=$1`,
		cart.ID,
	)
	if err != nil {
		shared.Log(shared.ERROR, err.Error())
		return nil, ErrInDatabase
	}

	for rows.Next() {
		var item CartItem
		err := rows.Scan(&item.ItemID, &item.Quantity)
		if err != nil {
			shared.Log(shared.ERROR, err.Error())
			return nil, ErrInDatabase
		}
		cart.Items = append(cart.Items, item)
	}

	return &cart, nil
}

func (r *PostgresRepository) Save(ctx context.Context, cart *Cart) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		shared.Log(shared.ERROR, err.Error())
		return ErrInDatabase
	}
	defer tx.Rollback(ctx)

	// 1. Lock cart row (REAL race protection)
	_, err = tx.Exec(ctx, `
		SELECT id 
		FROM carts 
		WHERE id = $1 
		FOR UPDATE
	`, cart.ID)

	if err != nil {
		shared.Log(shared.ERROR, err.Error())
		return ErrInDatabase
	}

	// 2. Update cart timestamp
	_, err = tx.Exec(ctx, `
		UPDATE carts 
		SET updated_at = now()
		WHERE id = $1
	`, cart.ID)

	if err != nil {
		shared.Log(shared.ERROR, err.Error())
		return ErrInDatabase
	}

	// 3. Delete old items (safe because locked)
	_, err = tx.Exec(ctx, `
		DELETE FROM cart_items 
		WHERE cart_id = $1
	`, cart.ID)

	if err != nil {
		shared.Log(shared.ERROR, err.Error())
		return ErrInDatabase
	}

	// 4. Insert fresh snapshot
	for _, item := range cart.Items {
		_, err = tx.Exec(ctx, `
			INSERT INTO cart_items (cart_id, item_id, quantity)
			VALUES ($1, $2, $3)
		`,
			cart.ID,
			item.ItemID,
			item.Quantity,
		)

		if err != nil {
			shared.Log(shared.ERROR, err.Error())
			return ErrInDatabase
		}
	}

	// 5. Commit atomically
	if err := tx.Commit(ctx); err != nil {
		shared.Log(shared.ERROR, err.Error())
		return ErrInDatabase
	}

	return nil
}

func (r *PostgresRepository) Empty(ctx context.Context, userID string) error {
	_, err := r.db.Exec(ctx,
		`DELETE FROM cart_items WHERE cart_id IN (SELECT id FROM carts WHERE user_id=$1)`,
		userID,
	)
	if err != nil {
		shared.Log(shared.ERROR, err.Error())
		return ErrInDatabase
	}

	_, err = r.db.Exec(ctx,
		`DELETE FROM carts WHERE user_id=$1`,
		userID,
	)
	if err != nil {
		shared.Log(shared.ERROR, err.Error())
		return ErrInDatabase
	}

	return nil
}
