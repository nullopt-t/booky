package user

import (
	"time"

	"github.com/google/uuid"
)

type UserCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateUserRequest struct {
	UserCredentials
}

type RegisterUserRequest struct {
	UserCredentials
}

type RegisterUserResponse struct {
	Email        string `json:"email,omitempty"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type CreateUserResponse struct {
	ID         uuid.UUID `json:"id"`
	Email      string    `json:"email"`
	IsInactive bool      `json:"is_inactive"`
	CreatedAt  string    `json:"created_at"`
	UpdatedAt  string    `json:"updated_at"`
}

type UserURIParam struct {
	ID uuid.UUID `json:"id"`
}

type UpdateUserRequest struct {
	Email      *string `json:"email"`
	IsInactive *bool   `json:"is_inactive"`
}

type UserResponse struct {
	ID         uuid.UUID `json:"id"`
	Email      string    `json:"email"`
	IsInactive bool      `json:"is_inactive"`
	DeletedAt  time.Time `json:"deleted_at,omitzero"`
	CreatedAt  time.Time `json:"created_at,omitzero"`
	UpdatedAt  time.Time `json:"updated_at,omitzero"`
}

type GetAllUsersResponse struct {
	Users []UserResponse `json:"users"`
	Page  int            `json:"page"`
	Limit int            `json:"limit"`
	Total int            `json:"total"`
}

type GetUserByEmailRequest struct {
	Email string `json:"email"`
}

type RefreshTokenRequest struct {
	Refresh_token string `json:"refresh_token"`
}

type RefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
}

type ForgetPasswordRequest struct {
	Email string `json:"email"`
}

type VerifyResetTokenRequest struct {
	Token       string `json:"token"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type VerifyEmailOTPRequest struct {
	Otp string `json:"otp"`
}
