package redis_client

import (
	"context"
	"fmt"
	"time"

	"github.com/phuslu/log"
	"github.com/redis/go-redis/v9"

	"github.com/bagasunix/gosnix/internal/domain/repository"
	"github.com/bagasunix/gosnix/pkg/errors"
)

type tokenCache struct {
	client *redis.Client
	logger *log.Logger
	prefix string
}

// Delete implements repository.TokenCacheRepository.
func (t *tokenCache) Delete(ctx context.Context, keys ...any) error {
	cacheKey := t.buildKey(keys...)

	if err := t.client.Del(ctx, cacheKey).Err(); err != nil {
		return errors.LogAndReturnError(t.logger, err, "failed to delete redis cache", "key", cacheKey)
	}
	return nil
}

// Get implements repository.TokenCacheRepository.
func (t *tokenCache) Get(ctx context.Context, keys ...any) (result string, err error) {
	cacheKey := t.buildKey(keys...)

	val, err := t.client.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		// Not found di cache bukan berarti error, return nil
		return "", nil
	} else if err != nil {
		return "", errors.LogAndReturnError(t.logger, err, "failed to get redis cache", "key", cacheKey)
	}

	return val, nil
}

// Set implements repository.TokenCacheRepository.
func (t *tokenCache) Set(ctx context.Context, ttl time.Duration, data any, keys ...any) error {
	cacheKey := t.buildKey(keys...)

	if err := t.client.Set(ctx, cacheKey, data, ttl).Err(); err != nil {
		return errors.LogAndReturnError(t.logger, err, "failed to set redis cache", "key", cacheKey)
	}
	return nil
}

// helper build key
func (t *tokenCache) buildKey(parts ...any) string {
	key := t.prefix
	for i, p := range parts {
		if i > 0 {
			key += ":"
		}
		key += fmt.Sprintf("%v", p)
	}
	return key
}

// NewTokenCache dengan prefix custom
func NewTokenCache(logger *log.Logger, client *redis.Client) repository.TokenCacheRepository {
	return &tokenCache{
		client: client,
		logger: logger,
		prefix: "tokens:",
	}
}
