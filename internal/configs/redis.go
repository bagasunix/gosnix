package configs

import (
	"context"
	"fmt"

	"github.com/phuslu/log"
	"github.com/redis/go-redis/v9"

	"github.com/bagasunix/gosnix/pkg/configs"
)

func InitRedis(ctx context.Context, logger *log.Logger, cfg *configs.Cfg) *redis.Client {
	options := &redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
	}

	client := redis.NewClient(options)

	// Test the connection
	_, err := client.Ping(ctx).Result()
	if err != nil {
		logger.Error().Err(err).Msg("Failed to connect to Redis")
		return nil
	}

	logger.Info().Msg("Connected to Redis")
	return client
}
