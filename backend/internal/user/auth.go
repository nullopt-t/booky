package user

import (
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
) error {
	_, err := s.userService.CreateUser(
		ctx,
		req.Email,
		req.Password,
	)
	if err != nil {
		return err
	}

	err = s.otpService.SendOTP(
		ctx,
		req.Email,
		"register",
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) VerifyOTP(
	ctx context.Context,
	req VerifyOTPRequest,
) error {
	err := s.otpService.VerifyOTP(
		ctx,
		req.Email,
		"register",
		req.Code,
	)
	if err != nil {
		return err
	}

	return s.userService.MarkEmailVerified(
		ctx,
		req.Email,
	)
}

func (s *AuthService) SendEmailOTP(
	ctx context.Context,
	email string,
) error {
	user, err := s.userService.GetUserByEmail(
		ctx,
		email,
	)
	if err != nil {
		return err
	}

	if user.EmailVerifiedAt != nil {
		return nil
	}

	return s.otpService.SendOTP(ctx,
		user.Email,
		"register",
	)
}
