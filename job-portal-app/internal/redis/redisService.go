// redisutil.go

package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/rs/zerolog/log"
	"job-portal-api/internal/models"
	"time"
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
			DB:       1,                // Default DB
		}),
	}
}

// SetData sets a key-value pair in Redis.
func (r *RedisClient) SetData(ctx context.Context, jid uint64, jobData models.Jobs) error {
	jobID := fmt.Sprint(jid)
	val, err := json.Marshal(jobData)
	fmt.Println("cache data unmarshaled:::", string(val))
	if err != nil {
		fmt.Println("cache set data unmrashal error:========", err)
		return err
	}
	err = r.client.Set(jobID, string(val), 10*time.Minute).Err()
	if err != nil {
		log.Error().Err(err).Msg("error while setting job data into cache")
	}

	return err
}

// GetData retrieves the value associated with a key from Redis.
func (r *RedisClient) GetData(ctx context.Context, jid uint64) (models.Jobs, error) {
	jobID := fmt.Sprint(jid)
	value, err := r.client.Get(jobID).Result()
	if err == redis.Nil {
		return models.Jobs{}, fmt.Errorf("job '%s' not found", jobID)
	}
	var job models.Jobs
	err = json.Unmarshal([]byte(value), &job)
	if err != nil {
		log.Error().Err(err).Msg("error while un-marshaling cache data")

		return models.Jobs{}, err
	}

	return job, nil
}

// Close closes the Redis client.
func (r *RedisClient) Close() {
	_ = r.client.Close()
}
