package redis

import (
	"github.com/phuslu/log"
	"github.com/redis/go-redis/v9"

	"github.com/bagasunix/gosnix/internal/domain/repository"
)

type RedisClient interface {
	GetHealth() repository.RedisRepository
}

type repo struct {
	health repository.RedisRepository
}

// GetHealth implements RedisClient.
func (r *repo) GetHealth() repository.RedisRepository {
	return r.health
}

func New(logger *log.Logger, redisClient *redis.Client) RedisClient {
	return &repo{health: NewHealthRepo(logger, redisClient)}
}
