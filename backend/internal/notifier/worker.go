package notifier

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

const LIST_KEY = "otp_email_queue"

type Worker struct {
	queue Queue
}

func NewNotifierWorker(queue Queue) *Worker {
	return &Worker{queue: queue}
}

func (w *Worker) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("context done, stopping worker")
			return
		default:
		}

		msg, err := w.queue.Dequeue(ctx, LIST_KEY)
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

		var payload OTPPayload
		if err := json.Unmarshal(msg.Payload, &payload); err != nil {
			continue
		}

		fmt.Printf("sending otp to %s\n", payload.Email)
	}
}
