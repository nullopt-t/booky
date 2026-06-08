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
const OTPMaxRetries = 3

type OTPPurpose string

const (
	OTPTypeLogin OTPPurpose = "login"
	OTPTypeReset OTPPurpose = "reset"
)

type Mailer interface {
	SendOTP(
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

type Service struct {
	otpRepo     OTPRepository
	userService UserService
	otpGen      OTPGenerator
	logger      log.Logger
	mailer      Mailer
}

func NewService(
	otpRepo OTPRepository,
	userService UserService,
	logger log.Logger,
	mailer Mailer,
) *Service {
	return &Service{
		otpRepo:     otpRepo,
		userService: userService,
		otpGen:      NewGenerator(),
		logger:      logger,
		mailer:      mailer,
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

func (s *Service) GenerateOTP(
	ctx context.Context,
	userID uuid.UUID,
	purpose string,
) (string, error) {
	otp, err := s.otpGen.GenerateOTP(purpose)
	if err != nil {
		return "", err
	}

	key := fmt.Sprintf("%s:%s",
		purpose,
		userID,
	)

	otpHash := jwt.Hash(otp)
	if err := s.otpRepo.Save(
		ctx,
		key,
		otpHash,
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
	user, err := s.userService.GetUserByID(
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

	if err := s.mailer.SendOTP(
		ctx,
		*user.Email,
		otp,
	); err != nil {
		return err
	}

	return nil
}

func (s *Service) VerifyOTP(
	ctx context.Context,
	userID uuid.UUID,
	purpose string,
	otp string,
) error {
	key := fmt.Sprintf("%s:%s",
		purpose,
		userID,
	)
	s.logger.Debug(
		"fetching otp hash....",
		log.Meta{"key": key},
	)
	value, err := s.otpRepo.Get(
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
			"key":   key,
			"otp":   otp,
			"value": value,
		},
	)

	if value != jwt.Hash(otp) {
		return invalidOTP()
	}
	s.logger.Debug(
		"deleting otp after successful verification",
		log.Meta{"key": key},
	)
	err = s.otpRepo.Delete(ctx, key)
	if err != nil {
		s.logger.Warn(
			"failed to delete otp",
			log.Meta{
				"key": key,
			},
		)
	}
	s.logger.Debug(
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
	user, err := s.userService.GetUserByID(
		ctx,
		userID,
	)
	if err != nil {
		return err
	}

	otp, err := s.otpGen.GenerateOTP(purpose)
	if err != nil {
		return err
	}

	optHash := jwt.Hash(otp)
	key := fmt.Sprintf("%s:%s",
		purpose,
		userID,
	)
	err = s.otpRepo.Save(
		ctx,
		key,
		optHash,
		OTPTTL,
	)
	if err != nil {
		return err
	}

	if err := s.mailer.SendOTP(ctx,
		*user.Email,
		otp); err != nil {
		return err
	}

	return nil
}
