package redis

import (
	"context"
	"fmt"
	"g6/blog-api/Delivery/bootstrap"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient interface {
	GetClient() *redis.Client
	Close() error
	Set(ctx context.Context, key string, value any, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
	Increment(ctx context.Context, key string) (int64, error)
	Decrement(ctx context.Context, key string) (int64, error)
	Expire(ctx context.Context, key string, expiration time.Duration) error
	GetCacheExpiry() time.Duration
}

type redisClient struct {
	client      *redis.Client
	cacheExpiry time.Duration
}

func NewRedisClient(env *bootstrap.Env) *redisClient {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", env.RedisHost, env.RedisPort),
		Password: env.RedisPassword,
		DB:       env.RedisDB,
	})

	return &redisClient{client: client, cacheExpiry: time.Duration(env.CacheExpirationSeconds) * time.Second}
}

func (r *redisClient) GetClient() *redis.Client {
	return r.client
}

func (r *redisClient) Close() error {
	if r.client != nil {
		return r.client.Close()
	}
	return nil
}

func (r *redisClient) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	if err := r.client.Set(ctx, key, value, expiration).Err(); err != nil {
		return fmt.Errorf("failed to set key %s: %w", key, err)
	}
	return nil
}

func (r *redisClient) Get(ctx context.Context, key string) (string, error) {
	value, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", fmt.Errorf("key %s does not exist", key)
		}
		return "", fmt.Errorf("failed to get key %s: %w", key, err)
	}
	return value, nil
}

func (r *redisClient) Delete(ctx context.Context, key string) error {
	if err := r.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete key %s: %w", key, err)
	}
	return nil
}

func (r *redisClient) Exists(ctx context.Context, key string) (bool, error) {
	exists, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check existence of key %s: %w", key, err)
	}
	return exists > 0, nil
}

func (r *redisClient) Increment(ctx context.Context, key string) (int64, error) {
	newValue, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to increment key %s: %w", key, err)
	}
	return newValue, nil
}

func (r *redisClient) Decrement(ctx context.Context, key string) (int64, error) {
	newValue, err := r.client.Decr(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to decrement key %s: %w", key, err)
	}
	return newValue, nil
}

func (r *redisClient) Expire(ctx context.Context, key string, expiration time.Duration) error {
	if err := r.client.Expire(ctx, key, expiration).Err(); err != nil {
		return fmt.Errorf("failed to set expiration for key %s: %w", key, err)
	}
	return nil
}

func (r *redisClient) GetCacheExpiry() time.Duration {
	if r.cacheExpiry <= 0 {
		return 1 * time.Hour
	}
	return r.cacheExpiry
}