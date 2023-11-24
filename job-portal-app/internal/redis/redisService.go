// redisutil.go

package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"job-portal-api/internal/config"
	"job-portal-api/internal/models"
	"time"

	"github.com/go-redis/redis"
	"github.com/rs/zerolog/log"
)

// RedisClient implements the Redis interface.
type RedisClient struct {
	client *redis.Client
}

// NewRedisClient initializes and returns a RedisClient.
func NewRedisClient(cfg config.Config) *RedisClient {
	return &RedisClient{
		client: redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", cfg.RedisConfig.Host, cfg.RedisConfig.Port), // Redis server address
			Password: cfg.RedisConfig.Password,                                         // No password
			DB:       cfg.RedisConfig.DB,                                               // Default DB
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

func (r *RedisClient) SaveOTPInCache(email, otp string) error {
	// Set the OTP in Redis with a specific expiration time (e.g., 5 minutes)
	//code, _ := json.Marshal(otp)

	err := r.client.Set(email, otp, 5*time.Minute).Err()
	if err != nil {
		return err
	}

	return nil
}
