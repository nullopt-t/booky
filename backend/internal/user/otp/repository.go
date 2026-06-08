package otp

import (
	"booky-backend/pkg/cache"
	"context"
	"time"
)

type OTPRepository interface {
	Save(ctx context.Context, key, value string, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
}

type Repository struct {
	cache cache.Cache
}

func NewOTPRepository(cache cache.Cache) *Repository {
	return &Repository{cache: cache}
}

func (r *Repository) Save(ctx context.Context, key, value string, expiration time.Duration) error {
	return r.cache.Set(ctx, key, value, expiration)
}

func (r *Repository) Get(ctx context.Context, key string) (string, error) {
	return r.cache.Get(ctx, key)
}

func (r *Repository) Delete(ctx context.Context, key string) error {
	return r.cache.Delete(ctx, key)
}
