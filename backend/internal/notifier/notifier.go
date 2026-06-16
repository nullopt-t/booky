package notifier

import (
	"context"
)

type MessageType string

const (
	MessageTypeOTP           MessageType = "email_otp"
	MessageTypeWelcome       MessageType = "email_welcome"
	MessageTypeResetPassword MessageType = "email_reset_password"
)

type OTPPayload struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

type ResetPasswordPayload struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

type WelcomePayload struct {
	Email string `json:"email"`
}

type Queue interface {
	Enqueue(ctx context.Context, msg Message) error
	Dequeue(ctx context.Context) (Message, error)
}
