package notifier

import (
	"booky-backend/pkg/log"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type MessageType string

const (
	MessageTypeOTP MessageType = "email_otp"
)

type Message struct {
	ID         uuid.UUID       `json:"id"`
	Type       MessageType     `json:"type"`
	Status     string          `json:"status"`
	Attempts   int             `json:"attempts"`
	Payload    json.RawMessage `json:"payload"`
	EnqueuedAt time.Time       `json:"enqueued_at"`
}

type OTPPayload struct {
	Email    string `json:"email"`
	CodeHash string `json:"code_hash"`
}

type Queue interface {
	Enqueue(ctx context.Context, key string, msg Message) error
	Dequeue(ctx context.Context, key string) (Message, error)
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

func (n *Notifier) SendOTP(ctx context.Context, email, otp string) error {
	payload := OTPPayload{
		Email:    email,
		CodeHash: otp,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}
	msg := Message{
		ID:         uuid.New(),
		Type:       MessageTypeOTP,
		Status:     "pending",
		Attempts:   0,
		Payload:    json.RawMessage(payloadBytes),
		EnqueuedAt: time.Now(),
	}
	n.logger.Info("enqueuing OTP message", log.Meta{
		"email": email,
		"otp":   otp,
	})
	err = n.queue.Enqueue(ctx, LIST_KEY, msg)
	if err != nil {
		return fmt.Errorf("failed to enqueue message: %w", err)
	}
	return nil
}
