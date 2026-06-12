package notifier

import (
	"booky-backend/pkg/log"
	"context"
	"encoding/json"
	"fmt"
	"time"
)

const LIST_KEY = "otp_email_queue"

type Worker struct {
	queue    Queue
	logger   log.Logger
	mailer   Mailer
	renderer *Renderer
}

func NewNotifierWorker(
	queue Queue,
	logger log.Logger,
	mailer Mailer,
	renderer *Renderer,
) *Worker {
	return &Worker{
		queue:    queue,
		logger:   logger,
		mailer:   mailer,
		renderer: renderer,
	}
}

func (w *Worker) handleMessage(
	msg Message,
) error {
	switch msg.Type {
	case MessageTypeOTP:
		var payload OTPPayload
		if err := json.Unmarshal(msg.Payload, &payload); err != nil {
			return err
		}
		tmpl, err := w.renderer.Render("otp", map[string]any{
			"Code": payload.Code,
		})
		if err != nil {
			return err
		}

		err = w.mailer.SendHTML([]string{payload.Email}, "OTP Code", tmpl)
		if err != nil {
			return err
		}
		fmt.Printf("otp %s sent to %s\n", payload.Code, payload.Email)

	case MessageTypeWelcome:
		var payload WelcomePayload
		if err := json.Unmarshal(msg.Payload, &payload); err != nil {
			return err
		}
		fmt.Printf("sending welcome to %s\n", payload.Email)
	default:
		return fmt.Errorf("unknown message type: %s", msg.Type)
	}
	return nil
}

func (w *Worker) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("context done, stopping worker")
			return
		default:
		}

		msg, err := w.queue.Dequeue(ctx)
		if err != nil {
			time.Sleep(500 * time.Millisecond) // basic backoff
			continue
		}

		fmt.Println("dequeued message:", map[string]any{
			"id":          msg.ID,
			"type":        msg.Type,
			"status":      msg.Status,
			"attempts":    msg.Attempts,
			"enqueued_at": msg.EnqueuedAt,
		})

		w.handleMessage(msg)
	}
}
