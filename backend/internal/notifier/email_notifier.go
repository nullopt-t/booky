package notifier

import (
	"booky-backend/pkg/log"
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type EmailNotifier struct {
	queue  Queue
	logger log.Logger
}

func NewEmailNotifier(
	queue Queue,
	logger log.Logger,
) *EmailNotifier {
	return &EmailNotifier{
		queue:  queue,
		logger: logger,
	}
}

func (n *EmailNotifier) NotifyOTP(
	ctx context.Context,
	email, otp string,
) error {
	n.logger.Debug("notify otp",
		log.Meta{
			"email": email,
			"otp":   otp,
		},
	)

	payload, err := json.Marshal(
		OTPPayload{
			Email: email,
			Code:  otp,
		},
	)
	if err != nil {
		return err
	}

	msg := &Message{
		ID:         uuid.New(),
		Type:       MessageTypeEmailOTP,
		Status:     "pending",
		Attempts:   0,
		Payload:    payload,
		EnqueuedAt: time.Now(),
	}
	return n.queue.Enqueue(ctx, msg)
}

func (n *EmailNotifier) NotifyWelcome(
	ctx context.Context,
	email string,
) error {
	payload, err := json.Marshal(
		WelcomePayload{
			Email: email,
		},
	)
	if err != nil {
		return err
	}
	msg := &Message{
		ID:         uuid.New(),
		Type:       MessageTypeEmailWelcome,
		Status:     "pending",
		Attempts:   0,
		Payload:    payload,
		EnqueuedAt: time.Now(),
	}
	return n.queue.Enqueue(ctx, msg)
}

func (n *EmailNotifier) NotifyResetPassword(
	ctx context.Context,
	email, token string,
) error {
	payload, err := json.Marshal(
		ResetPasswordPayload{
			Email: email,
			Token: token,
		},
	)
	if err != nil {
		return err
	}
	msg := &Message{
		ID:         uuid.New(),
		Type:       MessageTypeResetPassword,
		Status:     "pending",
		Attempts:   0,
		Payload:    payload,
		EnqueuedAt: time.Now(),
	}
	return n.queue.Enqueue(ctx, msg)
}
