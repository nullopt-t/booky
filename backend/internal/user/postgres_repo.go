package user

import (
	"booky-backend/internal/model"
	"booky-backend/pkg/api"
	"booky-backend/pkg/database"
	"context"

	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(ctx context.Context, qe database.QueryExecutor, email, passwordHashed string) (*model.User, error)
	GetUserByID(ctx context.Context, qe database.QueryExecutor, userID uuid.UUID) (*model.User, error)
	GetUserByEmail(ctx context.Context, qe database.QueryExecutor, email string) (*model.User, error)
	UpdateUser(ctx context.Context, qe database.QueryExecutor, userID uuid.UUID, user *model.User) error
	DeleteUser(ctx context.Context, qe database.QueryExecutor, userID uuid.UUID) error
	GetAllUsers(ctx context.Context, qe database.QueryExecutor, q *api.PageQuery) ([]*model.User, *api.Page, error)
}

type PostgresRepository struct {
}

func NewPostgresRepository() *PostgresRepository {
	return &PostgresRepository{}
}

func (r *PostgresRepository) CreateUser(ctx context.Context, qe database.QueryExecutor, email, passwordHashed string) (*model.User, error) {
	var createdUser model.User
	err := qe.QueryRow(ctx, "INSERT INTO users(email, password_hash) VALUES ($1, $2) RETURNING id, email", email, passwordHashed).Scan(&createdUser.ID, &createdUser.Email)
	if err != nil {
		return nil, database.MapError(err)
	}
	return &createdUser, nil
}

func (r *PostgresRepository) GetUserByID(ctx context.Context, qe database.QueryExecutor, userID uuid.UUID) (*model.User, error) {
	var existedUser model.User
	err := qe.QueryRow(ctx, "SELECT email, password_hash, is_inactive, failed_login_attempts, locked_until, created_at, updated_at, deleted_at FROM users WHERE id = $1 AND deleted_at IS NULL", userID).Scan(
		&existedUser.Email,
		&existedUser.PasswordHash,
		&existedUser.IsInactive,
		&existedUser.FailedLoginAttempts,
		&existedUser.LockedUntil,
		&existedUser.CreatedAt,
		&existedUser.UpdatedAt,
		&existedUser.DeletedAt)
	if err != nil {
		return nil, database.MapError(err)
	}
	return &existedUser, nil
}

func (r *PostgresRepository) GetUserByEmail(ctx context.Context, qe database.QueryExecutor, email string) (*model.User, error) {
	var existedUser model.User
	err := qe.QueryRow(ctx, "SELECT email, password_hash, is_inactive, failed_login_attempts, locked_until, created_at, updated_at, deleted_at FROM users WHERE email = $1 AND deleted_at IS NULL", email).Scan(
		&existedUser.Email,
		&existedUser.PasswordHash,
		&existedUser.IsInactive,
		&existedUser.FailedLoginAttempts,
		&existedUser.LockedUntil,
		&existedUser.CreatedAt,
		&existedUser.UpdatedAt,
		&existedUser.DeletedAt)
	if err != nil {
		return nil, database.MapError(err)
	}
	return &existedUser, nil
}

func (r *PostgresRepository) UpdateUser(ctx context.Context, qe database.QueryExecutor, userID uuid.UUID, user *model.User) error {
	_, err := qe.Exec(ctx, "UPDATE users SET email = $1, password_hash = $2, is_inactive = $3, failed_login_attempts = $4, locked_until = $5 WHERE id = $6 AND deleted_at IS NULL", user.Email, user.PasswordHash, user.IsInactive, user.FailedLoginAttempts, user.LockedUntil, userID)
	if err != nil {
		return database.MapError(err)
	}
	return nil
}

func (r *PostgresRepository) DeleteUser(ctx context.Context, qe database.QueryExecutor, userID uuid.UUID) error {
	_, err := qe.Exec(ctx, "UPDATE users SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL", userID)
	if err != nil {
		return database.MapError(err)
	}
	return nil
}

func (r *PostgresRepository) GetAllUsers(ctx context.Context, qe database.QueryExecutor, q *api.PageQuery) ([]*model.User, *api.Page, error) {
	offset := (q.Page - 1) * q.Limit
	var users []*model.User
	rows, err := qe.Query(ctx, "SELECT id, email, password_hash, is_inactive, failed_login_attempts, locked_until, created_at, updated_at, deleted_at FROM users WHERE deleted_at IS NULL LIMIT $1 OFFSET $2", q.Limit, offset)
	if err != nil {
		return nil, nil, database.MapError(err)
	}

	for rows.Next() {
		var user model.User
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.PasswordHash,
			&user.IsInactive,
			&user.FailedLoginAttempts,
			&user.LockedUntil,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.DeletedAt,
		)
		if err != nil {
			return nil, nil, database.MapError(err)
		}
		users = append(users, &user)
	}
	if err := rows.Err(); err != nil {
		return nil, nil, database.MapError(err)
	}

	var count int
	err = qe.QueryRow(ctx, "SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		return nil, nil, database.MapError(err)
	}

	page := &api.Page{
		Index: q.Page,
		Limit: q.Limit,
		Total: count,
	}

	return users, page, nil
}
