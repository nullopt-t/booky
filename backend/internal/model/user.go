package model

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	CustomerRole UserRole = "customer"
	AdminRole    UserRole = "admin"
	VendorRole   UserRole = "vendor"
)

type User struct {
	ID                  uuid.UUID
	Email               *string
	EmailOTP            *string
	EmailOTPExpiresAt   *time.Time
	EmailOTPAttempts    int
	IsEmailVerified     bool
	Phone               *string
	PhoneOTP            *string
	PhoneOTPAttempts    int
	PhoneOTPExpiresAt   *time.Time
	IsPhoneVerified     bool
	PasswordHash        string
	ResetToken          *string
	ResetTokenExpireAt  *time.Time
	LastResetRequestAt  *time.Time
	FailedResetAttempts int
	Role                UserRole
	FailedLoginAttempts int
	LockedUntil         *time.Time
	IsInactive          bool
	DeletedAt           *time.Time
	CreatedAt           time.Time
	UpdatedAt           time.Time
}
