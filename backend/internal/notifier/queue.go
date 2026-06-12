package notifier

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const QUEUE_KEY = "q:otp"

type Message struct {
	ID         uuid.UUID       `json:"id"`
	Type       MessageType     `json:"type"`
	Status     string          `json:"status"`
	Attempts   int             `json:"attempts"`
	Payload    json.RawMessage `json:"payload"`
	EnqueuedAt time.Time       `json:"enqueued_at"`
}

type redisJobQueue struct {
	client *redis.Client
}

func NewRedisJobQueue(client *redis.Client) *redisJobQueue {
	return &redisJobQueue{client: client}
}

func (q *redisJobQueue) Enqueue(ctx context.Context, msg Message) error {
	mm, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	err = q.client.LPush(ctx, QUEUE_KEY, mm).Err()
	if err != nil {
		return err
	}
	return nil
}

func (q *redisJobQueue) Dequeue(ctx context.Context) (Message, error) {
	var msg Message
	res, err := q.client.LPop(ctx, QUEUE_KEY).Result()
	if err != nil {
		return msg, err
	}
	if res != "" {
		err = json.Unmarshal([]byte(res), &msg)
		if err != nil {
			return msg, err
		}
	}
	return msg, nil
}
