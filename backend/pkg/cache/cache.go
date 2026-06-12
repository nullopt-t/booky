package cache

import (
	"context"
	"fmt"
	"net"
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

type MemoryCache struct {
	cache *cache.Cache
}

func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		cache: cache.New(
			cache.DefaultExpiration,
			cache.NoExpiration,
		),
	}
}

func (mc *MemoryCache) Get(_ context.Context, key string) (string, error) {
	val, ok := mc.cache.Get(key)
	if !ok {
		return "", fmt.Errorf("memory: key not found")
	}
	value, ok := val.(string)
	if !ok {
		return "", fmt.Errorf("memory: invalid value type")
	}
	return value, nil
}

func (mc *MemoryCache) Set(_ context.Context, key, value string, expiration time.Duration) error {
	err := mc.cache.Add(key, value, expiration)
	if err != nil {
		return fmt.Errorf("memory: set failed: %w", err)
	}
	return nil
}

func (mc *MemoryCache) Delete(_ context.Context, key string) error {
	mc.cache.Delete(key)
	return nil
}

func (mc *MemoryCache) Close(_ context.Context) error {
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
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", fmt.Errorf("redis: key not found")
		}
		if err := err.(net.Error); err != nil {
			return "", fmt.Errorf("redis: network error: %w", err)
		}
		return "", err
	}
	return val, nil
}

func (c *RedisCache) Set(
	ctx context.Context,
	key,
	value string,
	expiration time.Duration,
) error {
	if err := c.client.Set(ctx, key, value, expiration).Err(); err != nil {
		if err := err.(net.Error); err != nil {
			return fmt.Errorf("redis: network error: %w", err)
		}
		return err
	}
	return nil
}

func (c *RedisCache) Delete(
	ctx context.Context,
	key string,
) error {
	if err := c.client.Del(ctx, key).Err(); err != nil {
		if err := err.(net.Error); err != nil {
			return fmt.Errorf("redis: network error: %w", err)
		}
		return err
	}
	return nil
}

func (c *RedisCache) Close() error {
	err := c.client.Close()
	if err != nil {
		return fmt.Errorf("redis: close error: %w", err)
	}
	return nil
}
