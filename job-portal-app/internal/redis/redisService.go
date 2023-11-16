// redisutil.go

package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis"
)

// RedisClient implements the Redis interface.
type RedisClient struct {
	client *redis.Client
}

// NewRedisClient initializes and returns a RedisClient.
func NewRedisClient() *RedisClient {
	return &RedisClient{
		client: redis.NewClient(&redis.Options{
			Addr:     "localhost:6379", // Redis server address
			Password: "",               // No password
			DB:       0,                // Default DB
		}),
	}
}

// SetData sets a key-value pair in Redis.
func (r *RedisClient) SetData(ctx context.Context, key, value string) error {
	return r.client.Set(key, value, 0).Err()
}

// GetData retrieves the value associated with a key from Redis.
func (r *RedisClient) GetData(ctx context.Context, key string) (string, error) {
	value, err := r.client.Get(key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("key '%s' not found", key)
	} else if err != nil {
		return "", err
	}
	return value, nil
}

// Close closes the Redis client.
func (r *RedisClient) Close() {
	r.client.Close()
}
