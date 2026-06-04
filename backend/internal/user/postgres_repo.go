package user

import (
	"booky-backend/internal/model"
	"booky-backend/pkg/api"
	"booky-backend/pkg/database"
	"context"
	"time"

	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(ctx context.Context, qe database.QueryExecutor, email, passwordHashed string) (*model.User, error)
	GetUserByID(ctx context.Context, qe database.QueryExecutor, userID uuid.UUID) (*model.User, error)
	GetUserByEmail(ctx context.Context, qe database.QueryExecutor, email string) (*model.User, error)
	UpdateUser(ctx context.Context, qe database.QueryExecutor, userID uuid.UUID, user *model.User) error
	DeleteUser(ctx context.Context, qe database.QueryExecutor, userID uuid.UUID) error
	GetAllUsers(ctx context.Context, qe database.QueryExecutor, q *api.PageQuery) ([]*model.User, *api.Page, error)
	VerifyUserEmail(ctx context.Context, qe database.QueryExecutor, email string) error
	VerifyUserPhone(ctx context.Context, qe database.QueryExecutor, phone string) error

	// temporary
	SetUserEmailOTP(ctx context.Context, qe database.QueryExecutor, userID uuid.UUID, otp string, duration time.Duration) error
	ResetUserEmailOTP(ctx context.Context, qe database.QueryExecutor, userID uuid.UUID) error
	SetUserPhoneOTP(ctx context.Context, qe database.QueryExecutor, userID uuid.UUID, otp string, duration time.Duration) error
	ResetUserPhoneOTP(ctx context.Context, qe database.QueryExecutor, userID uuid.UUID) error

	// password
	UpdateUserPasswordHash(ctx context.Context, qe database.QueryExecutor, userID uuid.UUID, newPassHash string) error
	SetResetToken(ctx context.Context, qe database.QueryExecutor, userID uuid.UUID, token *string, duration time.Duration) error
	IncrementResetAttempts(ctx context.Context, qe database.QueryExecutor, userID uuid.UUID) error
	LockTokenResetFor(ctx context.Context, qe database.QueryExecutor, userID uuid.UUID, duration time.Duration) error
}

type PostgresRepository struct {
}

func NewPostgresRepository() *PostgresRepository {
	return &PostgresRepository{}
}

func (r *PostgresRepository) CreateUser(ctx context.Context, qe database.QueryExecutor, email, passwordHashed string) (*model.User, error) {
	var createdUser model.User
	err := qe.QueryRow(ctx, `INSERT INTO users(email, password_hash)
								 VALUES ($1, $2) RETURNING id, email`,
		email,
		passwordHashed).Scan(
		&createdUser.ID,
		&createdUser.Email,
	)
	return &createdUser, err
}

func (r *PostgresRepository) GetUserByID(ctx context.Context, qe database.QueryExecutor, userID uuid.UUID) (*model.User, error) {
	var existedUser model.User
	err := qe.QueryRow(ctx, `SELECT email, 
							  password_hash,
							  reset_token, 
							  is_inactive,
							  failed_login_attempts, 
						 	  locked_until,
						 	  created_at,
						 	  updated_at, 
						 	  deleted_at FROM users WHERE id = $1 AND deleted_at IS NULL`, userID).Scan(
		&existedUser.Email,
		&existedUser.PasswordHash,
		&existedUser.ResetToken,
		&existedUser.IsInactive,
		&existedUser.FailedLoginAttempts,
		&existedUser.LockedUntil,
		&existedUser.CreatedAt,
		&existedUser.UpdatedAt,
		&existedUser.DeletedAt)
	if err != nil {
		return nil, err
	}
	return &existedUser, nil
}

func (r *PostgresRepository) GetUserByEmail(ctx context.Context, qe database.QueryExecutor, email string) (*model.User, error) {
	var existedUser model.User
	err := qe.QueryRow(ctx,
		`SELECT id,
			email, 
			password_hash, 
			reset_token,
			reset_token_expires_at,
			is_inactive, 
			failed_login_attempts, 
			locked_until, 
			created_at, 
			updated_at, 
			deleted_at FROM users WHERE email = $1 AND deleted_at IS NULL`,
		email).Scan(
		&existedUser.ID,
		&existedUser.Email,
		&existedUser.PasswordHash,
		&existedUser.ResetToken,
		&existedUser.ResetTokenExpireAt,
		&existedUser.IsInactive,
		&existedUser.FailedLoginAttempts,
		&existedUser.LockedUntil,
		&existedUser.CreatedAt,
		&existedUser.UpdatedAt,
		&existedUser.DeletedAt)
	return &existedUser, err
}

func (r *PostgresRepository) UpdateUser(ctx context.Context, qe database.QueryExecutor, userID uuid.UUID, user *model.User) error {
	_, err := qe.Exec(ctx, `UPDATE users SET email = $1, 
								password_hash = $2, 
								is_inactive = $3,
								failed_login_attempts = $4, 
								locked_until = $5 WHERE id = $6 AND deleted_at IS NULL`,
		user.Email,
		user.PasswordHash,
		user.IsInactive,
		user.FailedLoginAttempts,
		user.LockedUntil,
		userID)
	return err
}

func (r *PostgresRepository) DeleteUser(ctx context.Context, qe database.QueryExecutor, userID uuid.UUID) error {
	_, err := qe.Exec(ctx, "UPDATE users SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL", userID)
	return err
}

func (r *PostgresRepository) GetAllUsers(ctx context.Context, qe database.QueryExecutor, q *api.PageQuery) ([]*model.User, *api.Page, error) {
	offset := (q.Page - 1) * q.Limit
	var users []*model.User
	rows, err := qe.Query(ctx, `SELECT id, 
									email, 
									password_hash,
									is_inactive,
									failed_login_attempts, 
									locked_until,
									created_at, 
									updated_at, 
									deleted_at FROM users WHERE deleted_at IS NULL LIMIT $1 OFFSET $2`, q.Limit, offset)
	if err != nil {
		return nil, nil, err
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
			return nil, nil, err
		}
		users = append(users, &user)
	}
	if err := rows.Err(); err != nil {
		return nil, nil, err
	}

	var count int
	err = qe.QueryRow(ctx, "SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		return nil, nil, err
	}

	page := &api.Page{
		Index: q.Page,
		Limit: q.Limit,
		Total: count,
	}

	return users, page, nil
}

func (r *PostgresRepository) SetUserEmailOTP(
	ctx context.Context,
	qe database.QueryExecutor,
	userID uuid.UUID,
	otp string,
	duration time.Duration) error {
	utc := time.Now().Add(duration)
	_, err := qe.Exec(ctx,
		`
		UPDATE users 
		SET email_otp = $2,
			email_otp_expires_at = $3
		WHERE id = $1 AND deleted_at IS NULL
		`, userID, otp, utc)
	return err
}

func (r *PostgresRepository) SetUserPhoneOTP(
	ctx context.Context,
	qe database.QueryExecutor,
	userID uuid.UUID,
	otp string,
	duration time.Duration) error {
	utc := time.Now().Add(duration)
	_, err := qe.Exec(ctx,
		`
		UPDATE users 
		SET phone_otp = $2,
			phone_otp_expires_at = $3
		WHERE id = $1 AND deleted_at IS NULL
		`, userID, otp, utc)
	return err
}

func (r *PostgresRepository) VerifyUserEmail(
	ctx context.Context,
	qe database.QueryExecutor,
	email string,
) error {
	_, err := qe.Exec(ctx,
		`
		UPDATE users 
		SET is_email_verified = TRUE
		WHERE email = $1 AND deleted_at IS NULL
		`, email)
	return err
}

func (r *PostgresRepository) ResetUserEmailOTP(
	ctx context.Context,
	qe database.QueryExecutor,
	userID uuid.UUID,
) error {
	_, err := qe.Exec(ctx,
		`
		UPDATE users 
		SET email_otp = NULL,
			email_otp_expires_at = NULL,
			email_otp_attempts = 0,
		WHERE id = $1 AND deleted_at IS NULL
		`, userID)
	return err
}

func (r *PostgresRepository) VerifyUserPhone(
	ctx context.Context,
	qe database.QueryExecutor,
	phone string,
) error {
	_, err := qe.Exec(ctx,
		`
		UPDATE users 
		SET is_email_verified = TRUE
		WHERE phone = $1 AND deleted_at IS NULL
		`, phone)
	return err
}

func (r *PostgresRepository) ResetUserPhoneOTP(
	ctx context.Context,
	qe database.QueryExecutor,
	userID uuid.UUID,
) error {
	_, err := qe.Exec(ctx,
		`
		UPDATE users 
		SET phone_otp = NULL,
			phone_otp_expires_at = NULL,
			phone_otp_attempts = 0,
		WHERE id = $1 AND deleted_at IS NULL
		`, userID)
	return err
}

func (r *PostgresRepository) UpdateUserPasswordHash(
	ctx context.Context,
	qe database.QueryExecutor,
	userID uuid.UUID,
	newPassHash string,
) error {
	_, err := qe.Exec(ctx,
		`
		UPDATE users 
		SET password_hash = $2 
		WHERE id = $1 AND deleted_at IS NULL
		`, userID, newPassHash)
	return err
}

func (r *PostgresRepository) SetResetToken(
	ctx context.Context,
	qe database.QueryExecutor,
	userID uuid.UUID,
	token *string,
	duration time.Duration,
) error {

	expireAt := time.Now().Add(duration)

	_, err := qe.Exec(ctx, `
        UPDATE users
        SET reset_token = $2,
            reset_token_expires_at = $3
        WHERE id = $1 AND deleted_at IS NULL
    `,
		userID,
		token,
		expireAt,
	)

	return err
}

func (r *PostgresRepository) IncrementResetAttempts(
	ctx context.Context,
	qe database.QueryExecutor,
	userID uuid.UUID,
) error {
	_, err := qe.Exec(ctx, `UPDATE users 
								SET failed_reset_attempts = failed_reset_attempts + 1
								WHERE id = $1 AND deleted_at IS NULL`,
		userID)
	return err
}

func (r *PostgresRepository) LockTokenResetFor(ctx context.Context, qe database.QueryExecutor, userID uuid.UUID, duration time.Duration) error {
	expireAt := time.Now().Add(duration)
	_, err := qe.Exec(ctx, `UPDATE users 
								SET reset_locked_until = $2,
								SET failed_reset_attempts = 0,
								WHERE id = $1 AND deleted_at IS NULL`,
		userID, expireAt)
	return err
}
