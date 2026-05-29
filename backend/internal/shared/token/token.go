package token

import (
	"booky-backend/internal/model"

	"github.com/google/uuid"
)

type UserSubject struct {
	UserID   uuid.UUID      `json:"user_id"`
	UserRole model.UserRole `json:"user_role"`
}
