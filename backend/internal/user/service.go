package user

import (
	"booky-backend/internal/model"
	"booky-backend/pkg/api"
	"booky-backend/pkg/database"
	"booky-backend/pkg/logger"
	"booky-backend/pkg/utils"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

const ResetTokenTTL = 15 * time.Minute

type UserService interface {
	CreateUser(ctx context.Context, user CreateUserRequest) (*model.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	GetAllUsers(ctx context.Context, q *api.PageQuery) ([]*model.User, *api.Page, error)
	UpdateUser(ctx context.Context, userID uuid.UUID, user *UpdateUserRequest) error
	DeleteUser(ctx context.Context, userID uuid.UUID) error

	ResendPhoneOTP(ctx context.Context, userID uuid.UUID) error
	VerifyPhoneOTP(ctx context.Context, userID uuid.UUID, otp string) error
	ResendEmailOTP(ctx context.Context, userID uuid.UUID) error
	VerifyEmailOTP(ctx context.Context, userID uuid.UUID, otp string) error
	CheckPassword(ctx context.Context, userID uuid.UUID, password string) error
	CheckPasswordResetToken(ctx context.Context, userID uuid.UUID, token string) error
	UpdatePassword(ctx context.Context, userID uuid.UUID, newPassword string) error
	SetResetToken(ctx context.Context, userID uuid.UUID, token *string) error
	IncrementResetAttempts(ctx context.Context, userID uuid.UUID) error
	LockTokenResetFor(ctx context.Context, userID uuid.UUID, duration time.Duration) error
}

type Service struct {
	dbExecuter database.Runner
	repo       UserRepository
}

func NewService(dbExecuter database.Runner, repo UserRepository) *Service {
	return &Service{
		dbExecuter: dbExecuter,
		repo:       repo,
	}
}

func (s *Service) CreateUser(ctx context.Context, user CreateUserRequest) (*model.User, error) {
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}
	var createdUser model.User
	err = s.dbExecuter.WithDB(ctx, func(db database.QueryExecutor) error {
		user, err := s.repo.CreateUser(ctx, db, user.Email, hashedPassword)
		if err != nil {
			return err
		}
		createdUser = *user
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &createdUser, nil
}

func (s *Service) GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	var existedUser model.User
	err := s.dbExecuter.WithDB(ctx, func(db database.QueryExecutor) error {
		user, err := s.repo.GetUserByID(ctx, db, id)
		if err != nil {
			return err
		}
		existedUser = *user
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &existedUser, nil
}

func (s *Service) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var existedUser model.User
	err := s.dbExecuter.WithDB(ctx, func(db database.QueryExecutor) error {
		user, err := s.repo.GetUserByEmail(ctx, db, email)
		if err != nil {
			return err
		}
		existedUser = *user
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &existedUser, nil
}

func (s *Service) GetAllUsers(ctx context.Context, q *api.PageQuery) ([]*model.User, *api.Page, error) {
	var users []*model.User
	var page *api.Page
	err := s.dbExecuter.WithDB(ctx, func(db database.QueryExecutor) error {
		var err error
		users, page, err = s.repo.GetAllUsers(ctx, db, q)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, nil, err
	}
	return users, page, nil
}

func (s *Service) UpdateUser(ctx context.Context, userID uuid.UUID, user *UpdateUserRequest) error {
	err := s.dbExecuter.WithDB(ctx, func(db database.QueryExecutor) error {
		existedUser, err := s.repo.GetUserByID(ctx, db, userID)
		if err != nil {
			return err
		}
		if user.Email != nil {
			existedUser.Email = *user.Email
		}
		if user.IsInactive != nil {
			existedUser.IsInactive = *user.IsInactive
		}
		return s.repo.UpdateUser(ctx, db, userID, existedUser)
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	err := s.dbExecuter.WithDB(ctx, func(db database.QueryExecutor) error {
		return s.repo.DeleteUser(ctx, db, userID)
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) ResendPhoneOTP(
	ctx context.Context,
	userID uuid.UUID,
) error {
	return s.dbExecuter.WithDB(ctx, func(db database.QueryExecutor) error {
		logger.Log(logger.DEBUG, "fetching user data...")
		user, err := s.repo.GetUserByID(
			ctx,
			db,
			userID,
		)
		if err != nil {
			return err
		}

		logger.Log(
			logger.DEBUG,
			"fetched user data",
			logger.LMeta{
				"user": user,
			},
		)

		if user.PhoneOTPExpiresAt != nil &&
			user.PhoneOTPExpiresAt.After(time.Now().Add(time.Minute*2)) {
			logger.Log(logger.DEBUG, "sending db otp")
			// send the saved otp
			return nil
		}

		logger.Log(logger.DEBUG, "generating otp...")
		otp, err := utils.GenerateOTP()
		if err != nil {
			return err
		}

		logger.Log(
			logger.DEBUG,
			"otp generated",
			logger.LMeta{
				"OTP": otp,
			},
		)

		err = s.repo.SetUserPhoneOTP(
			ctx,
			db,
			userID,
			otp,
			time.Minute*5,
		)
		if err != nil {
			return err
		}

		logger.Log(
			logger.DEBUG,
			"sending SMS",
		)

		// send sms with the new otp
		return nil
	})
}

func (s *Service) VerifyPhoneOTP(
	ctx context.Context,
	userID uuid.UUID,
	otp string,
) error {
	return s.dbExecuter.WithDB(ctx, func(db database.QueryExecutor) error {
		user, err := s.repo.GetUserByID(ctx, db, userID)
		if err != nil {
			return err
		}

		if user.IsEmailVerified {
			return fmt.Errorf("email already verified")
		}

		if strings.Compare(*user.EmailOTP, otp) != 0 {
			return fmt.Errorf("OTP does not match")
		}

		if err := s.repo.VerifyUserPhone(ctx, db, user.Email); err != nil {
			return err
		}
		return s.repo.ResetUserEmailOTP(ctx, db, user.ID)
	})
}

func (s *Service) ResendEmailOTP(
	ctx context.Context,
	userID uuid.UUID,
) error {
	return s.dbExecuter.WithDB(ctx, func(db database.QueryExecutor) error {
		logger.Log(logger.DEBUG, "fetching user data...")
		user, err := s.repo.GetUserByID(
			ctx,
			db,
			userID,
		)
		if err != nil {
			return err
		}

		logger.Log(
			logger.DEBUG,
			"fetched user data",
			logger.LMeta{
				"user": user,
			},
		)

		if user.EmailOTPExpiresAt != nil &&
			user.EmailOTPExpiresAt.After(time.Now().Add(time.Minute*2)) {
			logger.Log(logger.DEBUG, "sending db otp")
			// send the saved otp
			return nil
		}

		logger.Log(logger.DEBUG, "generating otp...")
		otp, err := utils.GenerateOTP()
		if err != nil {
			return err
		}

		logger.Log(
			logger.DEBUG,
			"otp generated",
			logger.LMeta{
				"OTP": otp,
			},
		)

		err = s.repo.SetUserEmailOTP(
			ctx,
			db,
			userID,
			otp,
			time.Minute*5,
		)
		if err != nil {
			return err
		}

		logger.Log(
			logger.DEBUG,
			"sending SMS",
		)

		// send email with the new otp
		return nil
	})
}

func (s *Service) VerifyEmailOTP(
	ctx context.Context,
	userID uuid.UUID,
	otp string,
) error {
	return s.dbExecuter.WithDB(ctx, func(db database.QueryExecutor) error {
		user, err := s.repo.GetUserByID(ctx, db, userID)
		if err != nil {
			return err
		}

		if user.IsEmailVerified {
			return fmt.Errorf("email already verified")
		}

		if strings.Compare(*user.EmailOTP, otp) != 0 {
			return fmt.Errorf("OTP does not match")
		}

		if err := s.repo.VerifyUserEmail(ctx, db, user.Email); err != nil {
			return err
		}
		return s.repo.ResetUserEmailOTP(ctx, db, user.ID)
	})
}

func (s *Service) CheckPassword(
	ctx context.Context,
	userID uuid.UUID,
	password string,
) error {
	return s.dbExecuter.WithDB(ctx, func(db database.QueryExecutor) error {
		user, err := s.repo.GetUserByID(ctx, db, userID)
		if err != nil {
			return err
		}

		if err := utils.ComparePassword(user.PasswordHash, password); err != nil {
			return fmt.Errorf("passord ain't matching")
		}

		return nil
	})
}

func (s *Service) CheckPasswordResetToken(
	ctx context.Context,
	userID uuid.UUID,
	token string,
) error {
	return s.dbExecuter.WithDB(ctx, func(db database.QueryExecutor) error {
		user, err := s.repo.GetUserByID(ctx, db, userID)
		if err != nil {
			return err
		}

		if strings.Compare(*user.ResetToken, token) != 0 {
			return fmt.Errorf("unknown token")
		}

		return nil
	})
}

func (s *Service) UpdatePassword(
	ctx context.Context,
	userID uuid.UUID,
	newPassword string,
) error {
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	return s.dbExecuter.WithDB(ctx, func(db database.QueryExecutor) error {
		return s.repo.UpdateUserPasswordHash(ctx, db, userID, hashedPassword)
	})
}

func (s *Service) SetResetToken(ctx context.Context, userID uuid.UUID, token *string) error {
	return s.dbExecuter.WithDB(ctx, func(db database.QueryExecutor) error {
		return s.repo.SetResetToken(ctx, db, userID, token, ResetTokenTTL)
	})
}

func (s *Service) IncrementResetAttempts(ctx context.Context, userID uuid.UUID) error {
	return s.dbExecuter.WithDB(ctx, func(db database.QueryExecutor) error {
		return s.repo.IncrementResetAttempts(ctx, db, userID)
	})
}

func (s *Service) LockTokenResetFor(ctx context.Context, userID uuid.UUID, duration time.Duration) error {
	return s.dbExecuter.WithDB(ctx, func(db database.QueryExecutor) error {
		return s.repo.LockTokenResetFor(ctx, db, userID, duration)
	})
}
