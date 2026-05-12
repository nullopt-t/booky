package cart

import (
	"booky-backend/internal/db"
	"booky-backend/internal/shared"
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type PostgresRepository struct{}

func NewPostgresRepository() CartRepository {
	return &PostgresRepository{}
}

func (r *PostgresRepository) Create(ctx context.Context, db db.DBQE, userID uuid.UUID) (*Cart, error) {
	var cart Cart
	err := db.QueryRow(ctx,
		`INSERT INTO carts (user_id, created_at, updated_at)
		 VALUES ($1, now(), now()) RETURNING id, created_at, updated_at`,
		userID,
	).Scan(&cart.ID, &cart.CreatedAt, &cart.UpdatedAt)

	if err != nil {
		shared.Log(shared.ERROR, err.Error())
		return nil, ErrInDatabase
	}

	return &cart, nil
}

func (r *PostgresRepository) GetByUserID(ctx context.Context, db db.DBQE, userID uuid.UUID) (*Cart, error) {
	var cart Cart
	err := db.QueryRow(ctx,
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

	rows, err := db.Query(ctx,
		`SELECT product_id, quantity FROM cart_items WHERE cart_id=$1`,
		cart.ID,
	)
	if err != nil {
		shared.Log(shared.ERROR, err.Error())
		return nil, ErrInDatabase
	}

	for rows.Next() {
		var item CartItem
		err := rows.Scan(&item.ProductID, &item.Quantity)
		if err != nil {
			shared.Log(shared.ERROR, err.Error())
			return nil, ErrInDatabase
		}
		cart.Items = append(cart.Items, item)
	}

	return &cart, nil
}


func (r *PostgresRepository) Save(ctx context.Context, db db.DBQE, cart *Cart) error {
	// 1. Lock cart row (REAL race protection)
	_, err := db.Exec(ctx, `
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
	_, err = db.Exec(ctx, `
		UPDATE carts 
		SET updated_at = now()
		WHERE id = $1
	`, cart.ID)

	if err != nil {
		shared.Log(shared.ERROR, err.Error())
		return ErrInDatabase
	}

	// 3. Delete old items (safe because locked)
	_, err = db.Exec(ctx, `
		DELETE FROM cart_items 
		WHERE cart_id = $1
	`, cart.ID)

	if err != nil {
		shared.Log(shared.ERROR, err.Error())
		return ErrInDatabase
	}

	// 4. Insert fresh snapshot
	for _, item := range cart.Items {
		_, err = db.Exec(ctx, `
			INSERT INTO cart_items (cart_id, product_id, quantity)
			VALUES ($1, $2, $3)
		`,
			cart.ID,
			item.ProductID,
			item.Quantity,
		)

		if err != nil {
			shared.Log(shared.ERROR, err.Error())
			return ErrInDatabase
		}
	}

	return nil
}

func (r *PostgresRepository) Empty(ctx context.Context, db db.DBQE, userID uuid.UUID) error {
	_, err := db.Exec(ctx,
		`DELETE FROM cart_items WHERE cart_id IN (SELECT id FROM carts WHERE user_id=$1)`,
		userID,
	)
	if err != nil {
		shared.Log(shared.ERROR, err.Error())
		return ErrInDatabase
	}

	_, err = db.Exec(ctx,
		`DELETE FROM carts WHERE user_id=$1`,
		userID,
	)
	if err != nil {
		shared.Log(shared.ERROR, err.Error())
		return ErrInDatabase
	}

	return nil
}

