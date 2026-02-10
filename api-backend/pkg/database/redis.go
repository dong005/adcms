package database

import (
	"adcms/internal/config"
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

var RDB *redis.Client

func InitRedis(cfg *config.RedisConfig) error {
	RDB = redis.NewClient(&redis.Options{
		Addr:     cfg.Addr(),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx := context.Background()
	if _, err := RDB.Ping(ctx).Result(); err != nil {
		return fmt.Errorf("failed to connect to redis: %w", err)
	}

	return nil
}

func CloseRedis() error {
	if RDB != nil {
		return RDB.Close()
	}
	return nil
}
