package db

import (
	"booky-backend/internal/config"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	pool *pgxpool.Pool
	cfg  *config.DatabaseConfig
}

func NewDatabase(cfg *config.DatabaseConfig) *DB {
	return &DB{cfg: cfg}
}

func (db *DB) Connect(ctx context.Context) error {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		db.cfg.DBHost,
		db.cfg.DBPort,
		db.cfg.DBUser,
		db.cfg.DBPassword,
		db.cfg.DBName,
	)

	var err error
	db.pool, err = pgxpool.New(ctx, dsn)
	if err != nil {
		return fmt.Errorf("db connection failed : %w", err)
	}

	err = db.pool.Ping(ctx)
	if err != nil {
		return fmt.Errorf("db connection failed : %w", err)

	}

	return nil
}

func (db *DB) Close() {
	db.pool.Close()
}

func (db *DB) GetPool() *pgxpool.Pool {
	return db.pool
}
