package database

import (
	"context"
	"fmt"
	"log"

	"go-trial/internal/config"

	"github.com/redis/go-redis/v9"
)

func NewRedis(cfg *config.RedisConfig) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
	})

	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Redis connected")
	return rdb
}