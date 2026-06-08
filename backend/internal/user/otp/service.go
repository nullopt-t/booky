package otp

import (
	"booky-backend/pkg/api/security"
	"booky-backend/pkg/log"
	"booky-backend/pkg/utils"
	"booky-backend/pkg/utils/jwt"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

const OTPLength = 6
const OTPTTL = 5 * time.Minute
const OTPMaxRetries = 3
const OTPKeyPrefix = "otp:"
const OTPAttemptsKeyPrefix = "otp_attempts:"

type OTPPurpose string

const (
	OTPTypeLogin OTPPurpose = "login"
	OTPTypeReset OTPPurpose = "reset"
)

type Mialer interface {
	SendOTP(userID uuid.UUID, otp string) error
}

type Service struct {
	otpRepo OTPRepository
	logger  log.Logger
	mailer  Mialer
}

func NewService(
	repo OTPRepository,
	logger log.Logger,
	mailer Mialer,
) *Service {
	return &Service{
		otpRepo: repo,
		logger:  logger,
		mailer:  mailer,
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

func (s *Service) generateOTP() (string, error) {
	otp, err := utils.GenerateOTP(OTPLength)
	if err != nil {
		return "", err
	}
	return otp, nil
}

func (s *Service) GenerateOTP(
	ctx context.Context,
	userID uuid.UUID,
	otpType OTPPurpose) (string, error) {
	otp, err := s.generateOTP()
	if err != nil {
		return "", err
	}
	key := fmt.Sprintf("%s:%s", otpType, userID)
	if err := s.otpRepo.Save(ctx, key, otp, OTPTTL); err != nil {
		return "", err
	}
	if err := s.mailer.SendOTP(userID, otp); err != nil {
		return "", err
	}
	return otp, nil
}

func (s *Service) VerifyOTP(
	ctx context.Context,
	userID uuid.UUID,
	otpType OTPPurpose,
	otp string,
) error {
	key := fmt.Sprintf("%s:%s", otpType, userID)
	s.logger.Debug(
		"fetching otp hash....",
		log.Meta{"key": key},
	)
	value, err := s.otpRepo.Get(ctx, key)
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
		log.Meta{"key": key},
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
		return err
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
	otpType OTPPurpose,
) error {
	otp, err := s.generateOTP()
	if err != nil {
		return err
	}
	optHash := jwt.Hash(otp)
	key := fmt.Sprintf("%s:%s", otpType, userID)
	err = s.otpRepo.Save(ctx, key, optHash, OTPTTL)
	if err != nil {
		return err
	}
	if err := s.mailer.SendOTP(userID, otp); err != nil {
		return err
	}
	return nil
}

func (s *Service) RevokeOTP(
	ctx context.Context,
	userID uuid.UUID,
	otpType OTPPurpose,
) error {
	key := fmt.Sprintf("%s:%s", otpType, userID)
	err := s.otpRepo.Delete(ctx, key)
	if err != nil {
		return err
	}
	return nil
}
