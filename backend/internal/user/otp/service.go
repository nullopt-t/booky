package otp

import (
	"booky-backend/internal/model"
	"booky-backend/pkg/api/security"
	"booky-backend/pkg/log"
	"booky-backend/pkg/utils/jwt"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

const OTPTTL = 15 * time.Minute

type OTPPurpose string

const (
	OTPTypeLogin OTPPurpose = "login"
	OTPTypeReset OTPPurpose = "reset"
)

type Notifier interface {
	NotifyOTP(
		ctx context.Context,
		to,
		otp string,
	) error
}

type Sender interface {
	SendSMS(
		ctx context.Context,
		to,
		otp string,
	) error
}

type UserService interface {
	GetUserByID(
		ctx context.Context,
		id uuid.UUID,
	) (*model.User, error)
}

type Generator interface {
	GenerateOTP(
		purpose string,
	) (string, error)
}

type Store interface {
	Save(
		ctx context.Context,
		key string,
		otp OTP,
		ttl time.Duration,
	) error

	Get(
		ctx context.Context,
		key string,
	) (*OTP, error)

	Increment(
		ctx context.Context,
		key string,
	) error

	Delete(
		ctx context.Context,
		key string,
	) error
}

type RateLimiter interface {
	AllowOTP(ctx context.Context, userID uuid.UUID) (bool, error)
}

type Service struct {
	store    Store
	userSrv  UserService
	gen      Generator
	limiter  RateLimiter
	logger   log.Logger
	notifier Notifier
}

func NewService(
	store Store,
	gen Generator,
	limiter RateLimiter,
	userSrv UserService,
	logger log.Logger,
	notifier Notifier,
) *Service {
	return &Service{
		store:    store,
		userSrv:  userSrv,
		gen:      gen,
		limiter:  limiter,
		logger:   logger,
		notifier: notifier,
	}
}

func invalidOTP() error {
	return security.NewSecureError(
		http.StatusUnauthorized,
		"EXPIRED_OR_INVALID_OTP",
		"expired or invalid OTP",
		nil,
	)
}

func (s *Service) genKey(purpose string, userID uuid.UUID) string {
	return fmt.Sprintf("%s:%s",
		purpose,
		userID,
	)
}

func (s *Service) GenerateOTP(
	ctx context.Context,
	userID uuid.UUID,
	purpose string,
) (string, error) {
	otp, err := s.gen.GenerateOTP(purpose)
	if err != nil {
		return "", err
	}

	key := s.genKey(purpose, userID)
	o := OTP{
		CodeHash: jwt.Hash(otp),
		Attempts: 0,
	}

	if err := s.store.Save(
		ctx,
		key,
		o,
		OTPTTL,
	); err != nil {
		return "", err
	}
	return otp, nil
}

func (s *Service) SendOTP(
	ctx context.Context,
	userID uuid.UUID,
	purpose string,
) error {
	allowed, err := s.limiter.AllowOTP(ctx, userID)
	if err != nil {
		return err
	}
	if !allowed {
		return security.NewSecureError(
			http.StatusTooManyRequests,
			"RATE_LIMIT_EXCEEDED",
			"rate limit exceeded",
			nil,
		)
	}

	user, err := s.userSrv.GetUserByID(
		ctx,
		userID,
	)
	if err != nil {
		return err
	}

	otp, err := s.GenerateOTP(
		ctx,
		userID,
		purpose,
	)
	if err != nil {
		return err
	}

	switch purpose {
	case "email":
		err = s.notifier.NotifyOTP(
			ctx,
			*user.Email,
			otp,
		)
	default:
		return fmt.Errorf("unsupported purpose: %s", purpose)
	}

	return err
}

func (s *Service) incrementAttempts(ctx context.Context, key string) error {
	err := s.store.Increment(
		ctx,
		key,
	)
	return err
}

func (s *Service) VerifyOTP(
	ctx context.Context,
	userID uuid.UUID,
	purpose string,
	otp string,
) error {
	var key string
	key = fmt.Sprintf("%s:%s",
		purpose,
		userID,
	)

	s.logger.Debug(
		"fetching otp hash....",
		log.Meta{"key": key},
	)

	uo, err := s.store.Get(
		ctx,
		key,
	)

	if err != nil {
		s.logger.Error(
			"failed to fetch otp hash",
			log.Meta{"key": key},
		)
		return security.NewSecureError(
			http.StatusUnauthorized,
			"EXPIRED_OR_INVALID_OTP",
			"expired or invalid OTP",
			err,
		)
	}
	s.logger.Debug(
		"otp hash fetched successfully",
		log.Meta{"key": key},
	)

	s.logger.Debug(
		"comparing otp hash",
		log.Meta{
			"key": key,
		},
	)

	if uo.CodeHash != jwt.Hash(otp) {
		err = s.incrementAttempts(ctx, key)
		if err != nil {
			s.logger.Error(
				"failed to increment attempts",
				log.Meta{"key": key},
			)
		}
		s.logger.Error(
			"invalid otp",
			log.Meta{"key": key},
		)
		return invalidOTP()
	}

	s.logger.Debug(
		"deleting otp after successful verification",
		log.Meta{"key": key},
	)
	err = s.store.Delete(ctx, key)
	if err != nil {
		s.logger.Warn(
			"failed to delete otp",
			log.Meta{
				"key": key,
			},
		)
	}
	s.logger.Info(
		"otp deleted successfully",
		log.Meta{"key": key},
	)
	return nil
}

func (s *Service) ResendOTP(
	ctx context.Context,
	userID uuid.UUID,
	purpose string,
) error {
	allowed, err := s.limiter.AllowOTP(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to check rate limit: %w", err)
	}
	if !allowed {
		return fmt.Errorf("rate limit exceeded")
	}

	// user, err := s.userSrv.GetUserByID(
	// 	ctx,
	// 	userID,
	// )
	// if err != nil {
	// 	return err
	// }

	otp, err := s.gen.GenerateOTP(purpose)
	if err != nil {
		return err
	}

	optHash := jwt.Hash(otp)
	key := fmt.Sprintf("%s:%s",
		purpose,
		userID,
	)
	err = s.store.Save(
		ctx,
		key,
		OTP{
			CodeHash: optHash,
		},
		OTPTTL,
	)
	if err != nil {
		return err
	}

	s.logger.Info(
		"otp saved successfully",
		log.Meta{"key": key},
	)

	s.logger.Info(
		"Sending OTP",
		log.Meta{
			"purpose": purpose,
			"userID":  userID,
			"otp":     otp,
		},
	)

	// switch purpose {
	// case "email":
	// 	err = s.mailer.SendOTP(ctx,
	// 		*user.Email,
	// 		otp)
	// case "phone":
	// 	err = s.mailer.SendOTP(ctx,
	// 		*user.Phone,
	// 		otp)
	// default:
	// 	return fmt.Errorf("unsupported purpose: %s", purpose)
	// }
	return err
}
