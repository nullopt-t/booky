package cache

import (
	"context"
	"errors"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/redis/go-redis/v9"
)

type Cache interface {
	Get(
		ctx context.Context,
		key string,
	) (string, error)

	Set(
		ctx context.Context,
		key string,
		value string,
		expiration time.Duration,
	) error

	Delete(
		ctx context.Context,
		key string,
	) error

	Close(
		ctx context.Context,
	) error
}

var (
	ErrNotExist    = errors.New("pair doesn't exist")
	ErrInvalidType = errors.New("invalid value type ")
	ErrSetFailed   = errors.New("failed to set a pair")
	ErrDelFailed   = errors.New("failed to delete a pair")
)

type MemoryCache struct {
	cache *cache.Cache
}

func NewMemoryCache(
	defaultTTL,
	cleanupInterval time.Duration,
) *MemoryCache {
	return &MemoryCache{
		cache: cache.New(
			defaultTTL,
			cleanupInterval,
		),
	}
}

func (mc *MemoryCache) Get(_ context.Context, key string) (string, error) {
	val, ok := mc.cache.Get(key)
	if !ok {
		return "", ErrNotExist
	}
	value, ok := val.(string)
	if !ok {
		return "", ErrInvalidType
	}
	return value, nil
}

func (mc *MemoryCache) Set(_ context.Context, key, value string, expiration time.Duration) error {
	err := mc.cache.Add(key, value, expiration)
	if err != nil {
		return err
	}
	return nil
}

func (mc *MemoryCache) Delete(_ context.Context, key string) error {
	mc.cache.Delete(key)
	_, ok := mc.cache.Get(key)
	if !ok {
		return ErrDelFailed
	}
	return nil
}

func (mc *MemoryCache) Clear() error {
	mc.cache.Flush()
	return nil
}

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(
	opt *redis.Options,
) *RedisCache {
	return &RedisCache{
		client: redis.NewClient(opt),
	}
}

func (c *RedisCache) Get(
	ctx context.Context,
	key string,
) (string, error) {
	cmd := c.client.Get(ctx, key)
	if cmd.Err() != nil {
		return "", cmd.Err()
	}
	return cmd.Val(), nil
}

func (c *RedisCache) Set(
	ctx context.Context,
	key,
	value string,
	expiration time.Duration,
) error {
	cmd := c.client.Set(ctx, key, value, expiration)
	if cmd.Err() != nil {
		return cmd.Err()
	}
	return nil
}

func (c *RedisCache) Delete(
	ctx context.Context,
	key string,
) error {
	cmd := c.client.Del(ctx, key)
	if cmd.Err() != nil {
		return cmd.Err()
	}
	return nil
}

func (c *RedisCache) Close() error {
	return c.client.Close()
}
