package redis_client

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/phuslu/log"
	"github.com/redis/go-redis/v9"

	"github.com/bagasunix/gosnix/internal/domain/entities"
	"github.com/bagasunix/gosnix/internal/domain/repository"
	"github.com/bagasunix/gosnix/pkg/errors"
)

type customerCache struct {
	client *redis.Client
	logger *log.Logger
	prefix string
}

// DeleteByPattern implements repository.CustomerCacheRepository.
func (c *customerCache) DeleteByPattern(ctx context.Context, pattern string) error {
	iter := c.client.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		if err := c.client.Del(ctx, iter.Val()).Err(); err != nil {
			return errors.LogAndReturnError(c.logger, err, "failed to delete by pattern", "pattern", pattern)
		}
	}
	if err := iter.Err(); err != nil {
		return errors.LogAndReturnError(c.logger, err, "failed to iterate scan", "pattern", pattern)
	}
	return nil
}

// helper build key
func (c *customerCache) buildKey(parts ...any) string {
	key := c.prefix
	for i, p := range parts {
		if i > 0 {
			key += ":"
		}
		key += fmt.Sprintf("%v", p)
	}
	return key
}

func (c *customerCache) Set(ctx context.Context, ttl time.Duration, data any, keys ...any) error {
	cacheKey := c.buildKey(keys...)

	if err := c.client.Set(ctx, cacheKey, data, ttl).Err(); err != nil {
		return errors.LogAndReturnError(c.logger, err, "failed to set redis cache", "key", cacheKey)
	}
	return nil
}

func (c *customerCache) GetCount(ctx context.Context, keys ...any) (result int, err error) {
	cacheKey := c.buildKey(keys...)

	val, err := c.client.Get(ctx, cacheKey).Int()
	if err == redis.Nil {
		// Not found di cache bukan berarti error, return nil
		return 0, nil
	} else if err != nil {
		return 0, errors.LogAndReturnError(c.logger, err, "failed to get redis cache", "key", cacheKey)
	}

	return val, nil
}

func (c *customerCache) Get(ctx context.Context, keys ...any) (result *string, err error) {
	cacheKey := c.buildKey(keys...)

	val, err := c.client.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		// Not found di cache bukan berarti error, return nil
		return nil, nil
	} else if err != nil {
		return nil, errors.LogAndReturnError(c.logger, err, "failed to get redis cache", "key", cacheKey)
	}

	return &val, nil
}

func (c *customerCache) GetWithValue(ctx context.Context, keys ...any) (result *entities.Customer, err error) {
	cacheKey := c.buildKey(keys...)

	val, err := c.client.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		// Not found di cache bukan berarti error, return nil
		return nil, nil
	} else if err != nil {
		return nil, errors.LogAndReturnError(c.logger, err, "failed to get redis cache", "key", cacheKey)
	}

	if err := json.Unmarshal([]byte(val), result); err != nil {
		return nil, errors.LogAndReturnError(c.logger, err, "failed to unmarshal customer from cache", "key", cacheKey)
	}

	return result, nil
}

func (c *customerCache) Delete(ctx context.Context, keys ...any) error {
	cacheKey := c.buildKey(keys...)

	if err := c.client.Del(ctx, cacheKey).Err(); err != nil {
		return errors.LogAndReturnError(c.logger, err, "failed to delete redis cache", "key", cacheKey)
	}
	return nil
}

// NewCustomerCache dengan prefix custom
func NewCustomerCache(logger *log.Logger, client *redis.Client) repository.CustomerCacheRepository {
	return &customerCache{
		client: client,
		logger: logger,
		prefix: "customers:",
	}
}
