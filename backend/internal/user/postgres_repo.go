package user

import (
	"booky-backend/internal/model"
	"booky-backend/pkg/api"
	"booky-backend/pkg/database"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type CreateUserParams struct {
	Email        *string
	Phone        *string
	PasswordHash string
}

type OTPRepository interface {
	// temporary
	SetUserEmailOTP(
		ctx context.Context,
		qe database.QueryExecutor,
		userID uuid.UUID,
		otp string,
		duration time.Duration,
	) error

	ResetUserEmailOTP(
		ctx context.Context,
		qe database.QueryExecutor,
		userID uuid.UUID,
	) error

	SetUserPhoneOTP(
		ctx context.Context,
		qe database.QueryExecutor,
		userID uuid.UUID,
		otp string,
		duration time.Duration,
	) error

	ResetUserPhoneOTP(
		ctx context.Context,
		qe database.QueryExecutor,
		userID uuid.UUID,
	) error
}

type UserRepository interface {
	CreateUser(
		ctx context.Context,
		qe database.QueryExecutor,
		params CreateUserParams,
	) (uuid.UUID, error)

	GetUserByID(
		ctx context.Context,
		qe database.QueryExecutor,
		userID uuid.UUID,
	) (*model.User, error)

	GetUserByEmail(
		ctx context.Context,
		qe database.QueryExecutor,
		email string,
	) (*model.User, error)

	GetUserByPhone(
		ctx context.Context,
		qe database.QueryExecutor,
		phone string,
	) (*model.User, error)

	DeleteUserByID(
		ctx context.Context,
		qe database.QueryExecutor,
		userID uuid.UUID,
	) error

	GetAllUsers(
		ctx context.Context,
		qe database.QueryExecutor,
		q api.PageQuery,
	) ([]model.User, api.Page, error)

	VerifyUserEmail(
		ctx context.Context,
		qe database.QueryExecutor,
		email string,
	) error

	VerifyUserPhone(
		ctx context.Context,
		qe database.QueryExecutor,
		phone string,
	) error

	// password
	UpdateUserPasswordHash(
		ctx context.Context,
		qe database.QueryExecutor,
		userID uuid.UUID,
		newPassHash string,
	) error

	SetResetToken(
		ctx context.Context,
		qe database.QueryExecutor,
		userID uuid.UUID,
		token string,
		duration time.Duration,
	) error

	IncrementResetAttempts(
		ctx context.Context,
		qe database.QueryExecutor,
		userID uuid.UUID,
	) error

	LockTokenResetFor(
		ctx context.Context,
		qe database.QueryExecutor,
		userID uuid.UUID,
		duration time.Duration,
	) error

	OTPRepository
}

type PostgresRepository struct {
}

func NewPostgresRepository() *PostgresRepository {
	return &PostgresRepository{}
}

func (r *PostgresRepository) CreateUser(
	ctx context.Context,
	qe database.QueryExecutor,
	params CreateUserParams,
) (uuid.UUID, error) {
	var createdUserID uuid.UUID
	err := qe.QueryRow(
		ctx,
		`INSERT INTO users(email, phone, password_hash)
		 VALUES ($1, $2, $3) RETURNING id`,
		params.Email,
		params.Phone,
		params.PasswordHash,
	).Scan(
		&createdUserID,
	)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create user: %w", err)
	}
	return createdUserID, nil
}

func (r *PostgresRepository) GetUserByID(
	ctx context.Context,
	qe database.QueryExecutor,
	userID uuid.UUID,
) (*model.User, error) {
	var existedUser model.User
	err := qe.QueryRow(
		ctx,
		`SELECT
				id,
				email,
				email_otp,
				email_otp_expires_at,
				email_otp_attempts,
				is_email_verified,
				phone,
				phone_otp,
				phone_otp_expires_at,
				phone_otp_attempts,
				is_phone_verified,
				role,
				password_hash,
				reset_token,
				failed_reset_attempts,
				last_reset_request_at,
				is_inactive,
				failed_login_attempts,
				locked_until,
				created_at,
				updated_at
				FROM users
				WHERE id = $1 AND deleted_at IS NULL`,
		userID,
	).Scan(
		&existedUser.ID,
		&existedUser.Email,
		&existedUser.EmailOTP,
		&existedUser.EmailOTPExpiresAt,
		&existedUser.EmailOTPAttempts,
		&existedUser.IsEmailVerified,
		&existedUser.Phone,
		&existedUser.PhoneOTP,
		&existedUser.PhoneOTPExpiresAt,
		&existedUser.PhoneOTPAttempts,
		&existedUser.IsPhoneVerified,
		&existedUser.Role,
		&existedUser.PasswordHash,
		&existedUser.ResetToken,
		&existedUser.FailedResetAttempts,
		&existedUser.LastResetRequestAt,
		&existedUser.IsInactive,
		&existedUser.FailedLoginAttempts,
		&existedUser.LockedUntil,
		&existedUser.CreatedAt,
		&existedUser.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	return &existedUser, nil
}

func (r *PostgresRepository) GetUserByEmail(
	ctx context.Context,
	qe database.QueryExecutor,
	email string,
) (*model.User, error) {
	var existedUser model.User
	err := qe.QueryRow(
		ctx,
		`SELECT
				id,
				email,
				email_otp,
				email_otp_expires_at,
				email_otp_attempts,
				is_email_verified,
				phone,
				phone_otp,
				phone_otp_expires_at,
				phone_otp_attempts,
				is_phone_verified,
				role,
				password_hash,
				reset_token,
				failed_reset_attempts,
				last_reset_request_at,
				is_inactive,
				failed_login_attempts,
				locked_until,
				created_at,
				updated_at
				FROM users
				WHERE email = $1 AND deleted_at IS NULL`,
		email,
	).Scan(
		&existedUser.ID,
		&existedUser.Email,
		&existedUser.EmailOTP,
		&existedUser.EmailOTPExpiresAt,
		&existedUser.EmailOTPAttempts,
		&existedUser.IsEmailVerified,
		&existedUser.Phone,
		&existedUser.PhoneOTP,
		&existedUser.PhoneOTPExpiresAt,
		&existedUser.PhoneOTPAttempts,
		&existedUser.IsPhoneVerified,
		&existedUser.Role,
		&existedUser.PasswordHash,
		&existedUser.ResetToken,
		&existedUser.FailedResetAttempts,
		&existedUser.LastResetRequestAt,
		&existedUser.IsInactive,
		&existedUser.FailedLoginAttempts,
		&existedUser.LockedUntil,
		&existedUser.CreatedAt,
		&existedUser.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return &existedUser, nil
}

func (r *PostgresRepository) GetUserByPhone(
	ctx context.Context,
	qe database.QueryExecutor,
	phone string,
) (*model.User, error) {
	var existedUser model.User
	err := qe.QueryRow(
		ctx,
		`SELECT
				id,
				email,
				email_otp,
				email_otp_expires_at,
				email_otp_attempts,
				is_email_verified,
				phone,
				phone_otp,
				phone_otp_expires_at,
				phone_otp_attempts,
				is_phone_verified,
				role,
				password_hash,
				reset_token,
				failed_reset_attempts,
				last_reset_request_at,
				is_inactive,
				failed_login_attempts,
				locked_until,
				created_at,
				updated_at
				FROM users
				WHERE phone = $1 AND deleted_at IS NULL`,
		phone,
	).Scan(
		&existedUser.ID,
		&existedUser.Email,
		&existedUser.EmailOTP,
		&existedUser.EmailOTPExpiresAt,
		&existedUser.EmailOTPAttempts,
		&existedUser.IsEmailVerified,
		&existedUser.Phone,
		&existedUser.PhoneOTP,
		&existedUser.PhoneOTPExpiresAt,
		&existedUser.PhoneOTPAttempts,
		&existedUser.IsPhoneVerified,
		&existedUser.Role,
		&existedUser.PasswordHash,
		&existedUser.ResetToken,
		&existedUser.FailedResetAttempts,
		&existedUser.LastResetRequestAt,
		&existedUser.IsInactive,
		&existedUser.FailedLoginAttempts,
		&existedUser.LockedUntil,
		&existedUser.CreatedAt,
		&existedUser.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by phone: %w", err)
	}
	return &existedUser, nil
}

func (r *PostgresRepository) DeleteUserByID(
	ctx context.Context,
	qe database.QueryExecutor,
	userID uuid.UUID,
) error {
	_, err := qe.Exec(
		ctx,
		`UPDATE users
		 		SET deleted_at = NOW()
		 		WHERE id = $1 AND deleted_at IS NULL`,
		userID,
	)
	if err != nil {
		return fmt.Errorf("failed to delete user by id: %w", err)
	}
	return nil
}

func (r *PostgresRepository) GetAllUsers(
	ctx context.Context,
	qe database.QueryExecutor,
	q api.PageQuery,
) ([]model.User, api.Page, error) {
	offset := (q.Page - 1) * q.PageSize
	var users []model.User
	var page api.Page
	rows, err := qe.Query(
		ctx,
		`SELECT
				id,
				email,
				email_otp,
				email_otp_expires_at,
				email_otp_attempts,
				is_email_verified,
				phone,
				phone_otp,
				phone_otp_expires_at,
				phone_otp_attempts,
				is_phone_verified,
				role,
				password_hash,
				reset_token,
				failed_reset_attempts,
				last_reset_request_at,
				is_inactive,
				failed_login_attempts,
				locked_until,
				created_at,
				updated_at
				FROM users
				WHERE deleted_at IS NULL LIMIT $1 OFFSET $2`,
		q.PageSize,
		offset,
	)
	if err != nil {
		return nil, page, fmt.Errorf("failed to get all users: %w", err)
	}

	for rows.Next() {
		var user model.User
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.EmailOTP,
			&user.EmailOTPExpiresAt,
			&user.EmailOTPAttempts,
			&user.IsEmailVerified,
			&user.Phone,
			&user.PhoneOTP,
			&user.PhoneOTPExpiresAt,
			&user.PhoneOTPAttempts,
			&user.IsPhoneVerified,
			&user.Role,
			&user.PasswordHash,
			&user.ResetToken,
			&user.FailedResetAttempts,
			&user.LastResetRequestAt,
			&user.IsInactive,
			&user.FailedLoginAttempts,
			&user.LockedUntil,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, page, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, page, fmt.Errorf("failed to iterate over users: %w", err)
	}

	var count int
	err = qe.QueryRow(ctx, "SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		return nil, page, fmt.Errorf("failed to count users: %w", err)
	}

	page = api.Page{
		Page:     q.Page,
		PageSize: q.PageSize,
		Total:    count,
	}

	return users, page, nil
}

func (r *PostgresRepository) SetUserEmailOTP(
	ctx context.Context,
	qe database.QueryExecutor,
	userID uuid.UUID,
	otp string,
	duration time.Duration,
) error {
	utc := time.Now().Add(duration)
	_, err := qe.Exec(ctx,
		`
		UPDATE users
		SET email_otp = $2,
			email_otp_expires_at = $3
		WHERE id = $1 AND deleted_at IS NULL
		`, userID, otp, utc)
	if err != nil {
		return fmt.Errorf("failed to set user email otp: %w", err)
	}
	return nil
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
	if err != nil {
		return fmt.Errorf("failed to set user phone otp: %w", err)
	}
	return nil
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
	if err != nil {
		return fmt.Errorf("failed to verify user email: %w", err)
	}
	return nil
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
			email_otp_attempts = 0
		WHERE id = $1 AND deleted_at IS NULL
		`, userID)
	if err != nil {
		return fmt.Errorf("failed to reset user email otp: %w", err)
	}
	return nil
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
	if err != nil {
		return fmt.Errorf("failed to verify user phone: %w", err)
	}
	return nil
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
			phone_otp_attempts = 0
		WHERE id = $1 AND deleted_at IS NULL
		`, userID)
	if err != nil {
		return fmt.Errorf("failed to reset user phone otp: %w", err)
	}
	return nil
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
	if err != nil {
		return fmt.Errorf("failed to update user password hash: %w", err)
	}
	return nil
}

func (r *PostgresRepository) SetResetToken(
	ctx context.Context,
	qe database.QueryExecutor,
	userID uuid.UUID,
	token string,
	duration time.Duration,
) error {

	expiresAt := time.Now().Add(duration)

	_, err := qe.Exec(ctx, `
        UPDATE users
        SET reset_token = $2,
            reset_token_expires_at = $3
        WHERE id = $1 AND deleted_at IS NULL
    `,
		userID,
		token,
		expiresAt,
	)

	if err != nil {
		return fmt.Errorf("failed to set reset token: %w", err)
	}
	return nil
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

	if err != nil {
		return fmt.Errorf("failed to increment reset attempts: %w", err)
	}
	return nil
}

func (r *PostgresRepository) LockTokenResetFor(
	ctx context.Context,
	qe database.QueryExecutor,
	userID uuid.UUID,
	duration time.Duration,
) error {
	expireAt := time.Now().Add(duration)
	_, err := qe.Exec(ctx, `UPDATE users
								SET reset_locked_until = $2,
								SET failed_reset_attempts = 0
								WHERE id = $1 AND deleted_at IS NULL`,
		userID, expireAt)
	if err != nil {
		return fmt.Errorf("failed to lock token reset: %w", err)
	}
	return nil
}
