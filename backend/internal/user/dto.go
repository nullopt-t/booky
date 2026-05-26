package user

import (
	"github.com/google/uuid"
)

type CreateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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
	CreatedAt  string    `json:"created_at"`
	UpdatedAt  string    `json:"updated_at"`
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
