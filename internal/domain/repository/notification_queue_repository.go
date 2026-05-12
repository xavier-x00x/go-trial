package repository

import (
	"context"
	"time"
)

type NotificationQueueRepository interface {
	Push(ctx context.Context, poID string) error
	Pop(ctx context.Context) (string, error)
	Length(ctx context.Context) (int, error)
	PushWithTTL(ctx context.Context, poID string, ttl time.Duration) error
	GetPending(ctx context.Context, limit int64) ([]string, error)
	Remove(ctx context.Context, poID string) error
	IsEmpty(ctx context.Context) (bool, error)
}