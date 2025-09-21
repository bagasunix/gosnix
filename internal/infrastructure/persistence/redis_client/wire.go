package redis_client

import (
	"github.com/phuslu/log"
	"github.com/redis/go-redis/v9"

	"github.com/bagasunix/gosnix/internal/domain/repository"
)

type RedisClient interface {
	GetHealth() repository.RedisRepository
	GetCustomerCache() repository.CustomerCacheRepository
}

type repo struct {
	health   repository.RedisRepository
	customer repository.CustomerCacheRepository
}

// GetCustiomerCache implements RedisClient.
func (r *repo) GetCustomerCache() repository.CustomerCacheRepository {
	return r.customer
}

// GetHealth implements RedisClient.
func (r *repo) GetHealth() repository.RedisRepository {
	return r.health
}

func New(logger *log.Logger, redisClient *redis.Client) RedisClient {
	return &repo{
		health:   NewHealthRepo(logger, redisClient),
		customer: NewCustomerCache(logger, redisClient),
	}
}
