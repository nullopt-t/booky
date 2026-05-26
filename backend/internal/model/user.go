package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                  uuid.UUID
	Email               string
	PasswordHash        string
	FailedLoginAttempts int
	LockedUntil         *time.Time
	IsInactive          bool
	DeletedAt           *time.Time
	CreatedAt           time.Time
	UpdatedAt           time.Time
}
