package redis

import (
	"context"
	"job-portal-api/internal/models"
)

// Redis interface defines the common Redis operations.
type Redis interface {
	SetData(ctx context.Context, jid uint64, jobData models.Jobs) error
	GetData(ctx context.Context, jid uint64) (string, error)
}
