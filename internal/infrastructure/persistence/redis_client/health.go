package redis_client

import (
	"context"
	"fmt"

	"github.com/phuslu/log"
	"github.com/redis/go-redis/v9"

	"github.com/bagasunix/gosnix/internal/domain/repository"
)

type redisProviderHealth struct {
	logger *log.Logger
	client *redis.Client
}

func NewHealthRepo(logger *log.Logger, client *redis.Client) repository.RedisRepository {
	g := new(redisProviderHealth)
	g.logger = logger
	g.client = client
	return g
}

func (g *redisProviderHealth) CheckRedis(ctx context.Context) error {
	if g == nil {
		return fmt.Errorf("redis health checker is not initialized")
	}

	if g.client == nil {
		return fmt.Errorf("redis client is not initialized")
	}

	return g.client.Ping(ctx).Err()
}
