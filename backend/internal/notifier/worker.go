package notifier

import (
	"booky-backend/pkg/config"
	"booky-backend/pkg/log"

	"context"
	"fmt"
	"time"
)

const LIST_KEY = "otp_email_queue"

type Mailer interface {
	SendHTML(to []string, subject, html string) error
}

type Worker struct {
	queue      Queue
	dispatcher *MessageDispatcher
	logger     log.Logger
	clientCfg  *config.ClientConfig
}

func NewNotifierWorker(
	queue Queue,
	dispatcher *MessageDispatcher,
	logger log.Logger,
	clientCfg *config.ClientConfig,
) *Worker {
	return &Worker{
		queue:      queue,
		dispatcher: dispatcher,
		logger:     logger,
		clientCfg:  clientCfg,
	}
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

		w.logger.Info(
			"dequeued message:",
			log.Meta{
				"id":          msg.ID,
				"type":        msg.Type,
				"status":      msg.Status,
				"attempts":    msg.Attempts,
				"enqueued_at": msg.EnqueuedAt,
			},
		)

		err = w.dispatcher.Dispatch(ctx, msg)
		if err != nil {
			w.logger.Error("failed to dispatch message", log.Meta{"error": err})
			msg.Status = "failed"
			msg.Attempts++
			err = w.queue.Enqueue(ctx, msg)
			if err != nil {
				w.logger.Error("failed to re-enqueue message", log.Meta{"error": err})
			}
			continue
		}
	}
}
