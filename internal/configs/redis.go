package configs

import (
	"context"
	"fmt"
	"time"

	"github.com/phuslu/log"
	"github.com/redis/go-redis/v9"

	"github.com/bagasunix/gosnix/pkg/configs"
	"github.com/bagasunix/gosnix/pkg/errors"
)

func InitRedis(ctx context.Context, logger *log.Logger, cfg *configs.Cfg) *redis.Client {
	options := &redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
		// Optional
		DialTimeout:  5 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		PoolSize:     10,
	}

	client := redis.NewClient(options)

	// Test the connection
	_, err := client.Ping(ctx).Result()
	errors.HandlerWithOSExit(logger, err, "init", "Redis", "config", options.Addr)

	logger.Info().Msg("Connected to Redis")
	return client
}
