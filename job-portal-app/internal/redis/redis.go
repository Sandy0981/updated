package redis

import "context"

// Redis interface defines the common Redis operations.
type Redis interface {
	SetData(ctx context.Context, key, value string) error
	GetData(ctx context.Context, key string) (string, error)
}
