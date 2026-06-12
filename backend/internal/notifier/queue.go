package notifier

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
)

type redisJobQueue struct {
	client *redis.Client
}

func NewRedisJobQueue(client *redis.Client) *redisJobQueue {
	return &redisJobQueue{client: client}
}

func (q *redisJobQueue) Enqueue(ctx context.Context, key string, msg Message) error {
	mm, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	err = q.client.LPush(ctx, key, mm).Err()
	if err != nil {
		return err
	}
	return nil
}

func (q *redisJobQueue) Dequeue(ctx context.Context, key string) (Message, error) {
	var msg Message
	res, err := q.client.LPop(ctx, key).Result()
	if err != nil {
		return Message{}, err
	}
	if res != "" {
		err = json.Unmarshal([]byte(res), &msg)
		if err != nil {
			return Message{}, err
		}
	}
	return msg, nil
}
