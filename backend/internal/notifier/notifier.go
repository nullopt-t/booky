package notifier

import (
	"booky-backend/pkg/log"
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type MessageType string

const (
	MessageTypeOTP     MessageType = "email_otp"
	MessageTypeWelcome MessageType = "email_welcome"
)

type OTPPayload struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

type WelcomePayload struct {
	Email string `json:"email"`
}

type Queue interface {
	Enqueue(ctx context.Context, msg Message) error
	Dequeue(ctx context.Context) (Message, error)
}

type Mailer interface {
	SendHTML(to []string, subject, html string) error
}

type Notifier struct {
	queue  Queue
	logger log.Logger
}

func NewNotifier(
	queue Queue,
	logger log.Logger,
) *Notifier {
	return &Notifier{
		queue:  queue,
		logger: logger,
	}
}

func (n *Notifier) NotifyOTP(
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
	return n.queue.Enqueue(ctx,
		Message{
			ID:         uuid.New(),
			Type:       MessageTypeOTP,
			Status:     "pending",
			Attempts:   0,
			Payload:    payload,
			EnqueuedAt: time.Now(),
		},
	)
}

func (n *Notifier) NotifyWelcome(
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
	return n.queue.Enqueue(ctx,
		Message{
			ID:         uuid.New(),
			Type:       MessageTypeWelcome,
			Status:     "pending",
			Attempts:   0,
			Payload:    payload,
			EnqueuedAt: time.Now(),
		},
	)
}
