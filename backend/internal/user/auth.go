package user

import (
	"booky-backend/internal/model"
	"booky-backend/pkg/log"
	"context"
)

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type AuthService struct {
	jwtService  *JwtService
	otpService  OTPService
	userService *UserService
	logger      log.Logger
}

func NewAuthService(
	logger log.Logger,
	userService *UserService,
	jwtService *JwtService,
	otpService OTPService,
) *AuthService {
	return &AuthService{
		jwtService:  jwtService,
		otpService:  otpService,
		userService: userService,
		logger:      logger,
	}
}

func (s *AuthService) Login(ctx context.Context, req LoginUserRequest) (*Tokens, error) {
	user, err := s.userService.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	accessToken, refreshToken, err := s.jwtService.GenerateTokenPair(user.ID.String(), string(user.Role))
	if err != nil {
		return nil, err
	}

	return &Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) Register(
	ctx context.Context,
	req RegisterUserRequest,
) (*model.User, error) {
	newUser, err := s.userService.Register(ctx, req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	// enqueue
	err = s.otpService.SendOTP(ctx, req.Email, "email")
	if err != nil {
		s.logger.Error(
			"failed to send OTP email",
			log.Meta{"user_id": newUser.ID.String()},
		)
	}

	return newUser, nil
}

func (s *AuthService) VerifyOTP(
	ctx context.Context,
	req VerifyOTPRequest,
) error {
	err := s.otpService.VerifyOTP(
		ctx,
		req.Purpose,
		req.Email,
		req.Code,
	)
	if err != nil {
		return err
	}

	err = s.userService.VerifyEmail(ctx, req.Email)
	if err != nil {
		return err
	}

	return nil
}
