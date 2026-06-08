package otp

import (
	"booky-backend/pkg/cache"
	"booky-backend/pkg/log"
	"context"
	"time"
)

type OTPRepository interface {
	Save(ctx context.Context, key, value string, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
}

type Repository struct {
	cache  cache.Cache
	logger log.Logger
}

func NewOTPRepository(
	cache cache.Cache,
	logger log.Logger,
) *Repository {
	return &Repository{cache: cache, logger: logger}
}

func (r *Repository) Save(
	ctx context.Context,
	key, value string,
	expiration time.Duration,
) error {
	r.logger.Debug(
		"saving otp",
		log.Meta{
			"key":        key,
			"expiration": expiration,
		},
	)
	return r.cache.Set(
		ctx,
		key,
		value,
		expiration,
	)
}

func (r *Repository) Get(
	ctx context.Context,
	key string,
) (string, error) {
	r.logger.Debug(
		"getting otp",
		log.Meta{
			"key": key,
		},
	)
	return r.cache.Get(
		ctx,
		key,
	)
}

func (r *Repository) Delete(
	ctx context.Context,
	key string,
) error {
	r.logger.Debug(
		"deleting otp",
		log.Meta{
			"key": key,
		},
	)
	return r.cache.Delete(
		ctx,
		key,
	)
}
