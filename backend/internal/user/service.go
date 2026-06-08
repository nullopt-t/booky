package user

import (
	"booky-backend/internal/model"
	"booky-backend/pkg/api"
	"booky-backend/pkg/api/security"
	"booky-backend/pkg/database"
	"booky-backend/pkg/log"
	"booky-backend/pkg/utils"
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

const ResetTokenTTL = 15 * time.Minute

type OTPService interface {
	SendOTP(
		ctx context.Context,
		userID uuid.UUID,
		purpose string,
	) error

	VerifyOTP(
		ctx context.Context,
		userID uuid.UUID,
		purpose string,
		otp string,
	) error
}

type UserService interface {
	CreateUser(
		ctx context.Context,
		user RegisterUserRequest,
	) (uuid.UUID, error)

	GetUserByID(
		ctx context.Context,
		id uuid.UUID,
	) (*model.User, error)

	GetUserByEmail(
		ctx context.Context,
		email string,
	) (*model.User, error)

	GetUserByPhone(
		ctx context.Context,
		phone string,
	) (*model.User, error)

	GetAllUsers(
		ctx context.Context,
		q api.PageQuery,
	) ([]model.User, api.Page, error)

	DeleteUserByID(
		ctx context.Context,
		userID uuid.UUID,
	) error

	CheckPassword(
		ctx context.Context,
		userID uuid.UUID,
		password string,
	) error

	CheckPasswordResetToken(
		ctx context.Context,
		userID uuid.UUID,
		token string,
	) error

	UpdatePassword(
		ctx context.Context,
		userID uuid.UUID,
		newPassword string,
	) error

	SetResetToken(
		ctx context.Context,
		userID uuid.UUID,
		token string,
	) error

	IncrementResetAttempts(
		ctx context.Context,
		userID uuid.UUID,
	) error

	LockTokenResetFor(
		ctx context.Context,
		userID uuid.UUID,
		duration time.Duration,
	) error

	VerifyUserEmail(
		ctx context.Context,
		userID uuid.UUID,
	) error
}

type Service struct {
	dbExecuter database.Runner
	repo       UserRepository
	logger     log.Logger
}

func NewService(
	dbExecuter database.Runner,
	repo UserRepository,
	logger log.Logger,
) *Service {
	return &Service{
		dbExecuter: dbExecuter,
		repo:       repo,
		logger:     logger,
	}
}

func (s *Service) CreateUser(
	ctx context.Context,
	req RegisterUserRequest,
) (uuid.UUID, error) {
	var createdID uuid.UUID

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return createdID, security.NewSecureError(
			http.StatusInternalServerError,
			security.CodeInternal,
			"failed to process user credentials",
			err,
		)
	}

	params := CreateUserParams{
		Email:        req.Email,
		Phone:        req.Phone,
		PasswordHash: hashedPassword,
	}

	err = s.dbExecuter.WithDB(
		ctx,
		func(db database.QueryExecutor) error {
			id, err := s.repo.CreateUser(
				ctx,
				db,
				params,
			)
			if err != nil {
				mappedErr := database.MapError(err)
				switch {
				case errors.Is(
					mappedErr,
					database.ErrConflict,
				):
					return security.NewSecureError(
						http.StatusConflict,
						security.CodeConflict,
						"user already exists",
						err,
					)
				default:
					return security.NewSecureError(
						http.StatusInternalServerError,
						security.CodeInternal,
						"failed to create a new user",
						err,
					)
				}
			}
			createdID = id
			return nil
		},
	)
	return createdID, err
}

func (s *Service) GetUserByID(
	ctx context.Context,
	id uuid.UUID,
) (*model.User, error) {
	var existedUser *model.User
	err := s.dbExecuter.WithDB(
		ctx,
		func(db database.QueryExecutor) error {
			user, err := s.repo.GetUserByID(
				ctx,
				db,
				id,
			)
			if err != nil {
				mappedErr := database.MapError(err)
				switch {
				case errors.Is(
					mappedErr,
					database.ErrNotFound,
				):
					return security.NewSecureError(
						http.StatusNotFound,
						security.CodeNotFound,
						err.Error(),
						err,
					)
				default:
					return security.NewSecureError(
						http.StatusInternalServerError,
						security.CodeInternal,
						"failed to fetch a user",
						err,
					)
				}
			}
			existedUser = user
			return nil
		},
	)
	return existedUser, err
}

func (s *Service) GetUserByEmail(
	ctx context.Context,
	email string,
) (*model.User, error) {
	var existedUser *model.User
	err := s.dbExecuter.WithDB(
		ctx,
		func(db database.QueryExecutor) error {
			user, err := s.repo.GetUserByEmail(
				ctx,
				db,
				email,
			)
			if err != nil {
				mappedErr := database.MapError(err)
				switch {
				case errors.Is(
					mappedErr,
					database.ErrNotFound,
				):
					return security.NewSecureError(
						http.StatusNotFound,
						security.CodeNotFound,
						err.Error(),
						err,
					)
				default:
					return security.NewSecureError(
						http.StatusInternalServerError,
						security.CodeInternal,
						"failed to fetch a user",
						err,
					)
				}
			}
			existedUser = user
			return nil
		},
	)
	return existedUser, err
}

func (s *Service) GetUserByPhone(
	ctx context.Context,
	phone string,
) (*model.User, error) {
	var existedUser model.User
	err := s.dbExecuter.WithDB(
		ctx,
		func(db database.QueryExecutor) error {
			user, err := s.repo.GetUserByPhone(
				ctx,
				db,
				phone,
			)
			if err != nil {
				mappedErr := database.MapError(err)
				switch {
				case errors.Is(
					mappedErr,
					database.ErrNotFound,
				):
					return security.NewSecureError(
						http.StatusNotFound,
						security.CodeNotFound,
						err.Error(),
						err,
					)
				default:
					return security.NewSecureError(
						http.StatusInternalServerError,
						security.CodeInternal,
						"failed to fetch a user",
						err,
					)
				}
			}
			existedUser = *user
			return nil
		},
	)
	if err != nil {
		return nil, err
	}
	return &existedUser, nil
}

func (s *Service) GetAllUsers(
	ctx context.Context,
	q api.PageQuery,
) ([]model.User, api.Page, error) {
	var users []model.User
	var page api.Page
	err := s.dbExecuter.WithDB(
		ctx,
		func(db database.QueryExecutor) error {
			var err error
			users, page, err = s.repo.GetAllUsers(
				ctx,
				db,
				q,
			)
			if err != nil {
				mappedErr := database.MapError(err)
				switch {
				case errors.Is(
					mappedErr,
					database.ErrNotFound,
				):
					return security.NewSecureError(
						http.StatusNotFound,
						security.CodeNotFound,
						err.Error(),
						err,
					)
				default:
					return security.NewSecureError(
						http.StatusInternalServerError,
						security.CodeInternal,
						"failed to fetch users",
						err,
					)
				}
			}
			return nil
		},
	)
	if err != nil {
		return nil, page, err
	}
	return users, page, nil
}

func (s *Service) DeleteUserByID(
	ctx context.Context,
	userID uuid.UUID,
) error {
	err := s.dbExecuter.WithDB(
		ctx,
		func(db database.QueryExecutor) error {
			return s.repo.DeleteUserByID(
				ctx,
				db,
				userID,
			)
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) CheckPassword(
	ctx context.Context,
	userID uuid.UUID,
	password string,
) error {
	return s.dbExecuter.WithDB(
		ctx,
		func(db database.QueryExecutor) error {
			user, err := s.repo.GetUserByID(
				ctx,
				db,
				userID,
			)
			if err != nil {
				return security.NewSecureError(
					http.StatusInternalServerError,
					security.CodeInternal,
					"failed to fetch a user",
					err,
				)
			}

			if err = utils.ComparePassword(
				user.PasswordHash,
				password,
			); err != nil {
				return security.NewSecureError(
					http.StatusUnauthorized,
					security.CodeUnauthorized,
					"idenifier or password is incorrect",
					err,
				)
			}

			return nil
		},
	)
}

func (s *Service) CheckPasswordResetToken(
	ctx context.Context,
	userID uuid.UUID,
	token string,
) error {
	return s.dbExecuter.WithDB(
		ctx,
		func(db database.QueryExecutor) error {
			user, err := s.repo.GetUserByID(
				ctx,
				db,
				userID,
			)
			if err != nil {
				return security.NewSecureError(
					http.StatusInternalServerError,
					security.CodeInternal,
					"failed to fetch a user",
					err,
				)
			}

			if strings.Compare(
				*user.ResetToken,
				token,
			) != 0 {
				return security.NewSecureError(
					http.StatusUnauthorized,
					security.CodeUnauthorized,
					"invalid token",
					errors.New("invalid token"),
				)
			}

			return nil
		},
	)
}

func (s *Service) UpdatePassword(
	ctx context.Context,
	userID uuid.UUID,
	newPassword string,
) error {
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return security.NewSecureError(
			http.StatusInternalServerError,
			security.CodeInternal,
			"failed to process the credentials",
			nil,
		)
	}

	err = s.dbExecuter.WithDB(
		ctx,
		func(db database.QueryExecutor) error {
			return s.repo.UpdateUserPasswordHash(
				ctx,
				db,
				userID,
				hashedPassword,
			)
		},
	)
	if err != nil {
		return security.NewSecureError(
			http.StatusInternalServerError,
			security.CodeInternal,
			"failed to update password",
			err,
		)
	}
	return nil
}

func (s *Service) SetResetToken(
	ctx context.Context,
	userID uuid.UUID,
	token string,
) error {
	err := s.dbExecuter.WithDB(
		ctx,
		func(db database.QueryExecutor) error {
			return s.repo.SetResetToken(
				ctx,
				db,
				userID,
				token,
				ResetTokenTTL,
			)
		},
	)
	if err != nil {
		return security.NewSecureError(
			http.StatusInternalServerError,
			security.CodeInternal,
			"failed to set reset token",
			err,
		)
	}
	return nil
}

func (s *Service) IncrementResetAttempts(
	ctx context.Context,
	userID uuid.UUID,
) error {
	err := s.dbExecuter.WithDB(
		ctx,
		func(db database.QueryExecutor) error {
			return s.repo.IncrementResetAttempts(
				ctx,
				db,
				userID,
			)
		},
	)
	if err != nil {
		return security.NewSecureError(
			http.StatusInternalServerError,
			security.CodeInternal,
			"failed to increment reset attempts",
			err,
		)
	}
	return nil
}

func (s *Service) LockTokenResetFor(
	ctx context.Context,
	userID uuid.UUID,
	duration time.Duration,
) error {
	err := s.dbExecuter.WithDB(
		ctx,
		func(db database.QueryExecutor) error {
			return s.repo.LockTokenResetFor(
				ctx,
				db,
				userID,
				duration,
			)
		})
	if err != nil {
		return security.NewSecureError(
			http.StatusInternalServerError,
			security.CodeInternal,
			"failed to lock token reset",
			err,
		)
	}
	return nil
}

func (s *Service) VerifyUserEmail(
	ctx context.Context,
	userID uuid.UUID,
) error {
	s.logger.Debug(
		"VerifyUserEmail",
		log.Meta{
			"context": ctx,
			"userID":  userID,
		},
	)
	err := s.dbExecuter.WithDB(
		ctx,
		func(db database.QueryExecutor) error {
			return s.repo.VerifyUserEmail(
				ctx,
				db,
				userID,
			)
		},
	)
	if err != nil {
		return security.NewSecureError(
			http.StatusInternalServerError,
			security.CodeInternal,
			"failed to verify user email",
			err,
		)
	}
	return nil
}
