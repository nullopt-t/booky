package otp

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const (
	MaxOTPRequests = 5
)

type Limiter struct {
	client *redis.Client
}

func NewRateLimiter(client *redis.Client) *Limiter {
	return &Limiter{client: client}
}

func (r *Limiter) AllowOTP(ctx context.Context, userID uuid.UUID) (bool, error) {
	key := fmt.Sprintf("otp:requests:%s", userID)

	now := time.Now().Unix()
	window := int64(time.Hour.Seconds())

	// Remove requests older than one hour
	if err := r.client.ZRemRangeByScore(
		ctx,
		key,
		"0",
		strconv.FormatInt(now-window, 10),
	).Err(); err != nil {
		return false, err
	}

	// Count remaining requests
	count, err := r.client.ZCard(ctx, key).Result()
	if err != nil {
		return false, err
	}

	if count >= MaxOTPRequests {
		return false, nil
	}

	// Store current request
	err = r.client.ZAdd(ctx, key, redis.Z{
		Score:  float64(now),
		Member: uuid.NewString(),
	}).Err()
	if err != nil {
		return false, err
	}

	// Auto cleanup
	r.client.Expire(ctx, key, time.Hour)
	return true, nil
}
