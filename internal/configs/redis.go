package configs

import (
	"context"

	"github.com/phuslu/log"
	"github.com/redis/go-redis/v9"

	"github.com/bagasunix/gosnix/pkg/configs"
	"github.com/bagasunix/gosnix/pkg/errors"
	"github.com/bagasunix/gosnix/pkg/utils"
)

func InitRedis(ctx context.Context, logger *log.Logger, cfg *configs.Cfg) *redis.Client {
	CfgBuild := &utils.DBConfig{
		Driver:       cfg.Redis.Driver,
		Host:         cfg.Redis.Host,
		Port:         cfg.Redis.Port,
		User:         cfg.Redis.User,
		Password:     cfg.Redis.Password,
		DatabaseName: cfg.Redis.DB,
	}
	return NewRedisDB(ctx, CfgBuild, logger)
}

// NewRedisDB initializes and returns a Redis client
func NewRedisDB(ctx context.Context, cfg *utils.DBConfig, logger *log.Logger) *redis.Client {
	dsn := cfg.GetDSN()

	opt, err := redis.ParseURL(dsn)
	if err != nil {
		errors.HandlerWithOSExit(logger, err, "init", "Redis", "parseURL", dsn)
	}

	client := redis.NewClient(opt)
	// Test the connection
	if _, err = client.Ping(ctx).Result(); err != nil {
		errors.HandlerWithOSExit(logger, err, "init", "Redis", "ping", dsn)
	}

	logger.Info().Msg("Connected to Redis")
	return client
}
