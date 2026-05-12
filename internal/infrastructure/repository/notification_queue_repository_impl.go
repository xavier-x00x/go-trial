package repository

import (
	"context"
	"fmt"
	"time"

	"go-trial/internal/domain/repository"

	"github.com/redis/go-redis/v9"
)

const (
	notificationQueueKey = "notification:po:queue"
)

type notificationQueueRepository struct {
	rdb *redis.Client
}

func NewNotificationQueueRepository(rdb *redis.Client) repository.NotificationQueueRepository {
	return &notificationQueueRepository{rdb: rdb}
}

func (r *notificationQueueRepository) Push(ctx context.Context, poID string) error {
	return r.rdb.RPush(ctx, notificationQueueKey, poID).Err()
}

func (r *notificationQueueRepository) Pop(ctx context.Context) (string, error) {
	result, err := r.rdb.LPop(ctx, notificationQueueKey).Result()
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("failed to pop from notification queue: %w", err)
	}
	return result, nil
}

func (r *notificationQueueRepository) Length(ctx context.Context) (int, error) {
	length, err := r.rdb.LLen(ctx, notificationQueueKey).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get queue length: %w", err)
	}
	return int(length), nil
}

func (r *notificationQueueRepository) PushWithTTL(ctx context.Context, poID string, ttl time.Duration) error {
	return r.rdb.RPush(ctx, notificationQueueKey, poID).Err()
}

func (r *notificationQueueRepository) GetPending(ctx context.Context, limit int64) ([]string, error) {
	results, err := r.rdb.LRange(ctx, notificationQueueKey, 0, limit-1).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get pending notifications: %w", err)
	}
	return results, nil
}

func (r *notificationQueueRepository) Remove(ctx context.Context, poID string) error {
	return r.rdb.LRem(ctx, notificationQueueKey, 0, poID).Err()
}

func (r *notificationQueueRepository) IsEmpty(ctx context.Context) (bool, error) {
	length, err := r.Length(ctx)
	if err != nil {
		return false, err
	}
	return length == 0, nil
}