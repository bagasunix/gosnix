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

type vehicleCache struct {
	client *redis.Client
	logger *log.Logger
	prefix string
}

// buildKey implements repository.VehicleCacheRepository.
func (v *vehicleCache) buildKey(parts ...any) string {
	key := v.prefix
	for i, p := range parts {
		if i > 0 {
			key += ":"
		}
		key += fmt.Sprintf("%v", p)
	}
	return key
}

// Delete implements repository.VehicleCacheRepository.
func (v *vehicleCache) Delete(ctx context.Context, keys ...any) error {
	cacheKey := v.buildKey(keys...)

	if err := v.client.Del(ctx, cacheKey).Err(); err != nil {
		return errors.LogAndReturnError(v.logger, err, "failed to delete redis cache", "key", cacheKey)
	}
	return nil
}

// DeleteByPattern implements repository.VehicleCacheRepository.
func (v *vehicleCache) DeleteByPattern(ctx context.Context, pattern string) error {
	iter := v.client.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		if err := v.client.Del(ctx, iter.Val()).Err(); err != nil {
			return errors.LogAndReturnError(v.logger, err, "failed to delete by pattern", "pattern", pattern)
		}
	}
	if err := iter.Err(); err != nil {
		return errors.LogAndReturnError(v.logger, err, "failed to iterate scan", "pattern", pattern)
	}
	return nil
}

// Get implements repository.VehicleCacheRepository.
func (v *vehicleCache) Get(ctx context.Context, keys ...any) (result *string, err error) {
	cacheKey := v.buildKey(keys...)

	val, err := v.client.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		// Not found di cache bukan berarti error, return nil
		return nil, nil
	} else if err != nil {
		return nil, errors.LogAndReturnError(v.logger, err, "failed to get redis cache", "key", cacheKey)
	}

	return &val, nil
}

// GetCount implements repository.VehicleCacheRepository.
func (v *vehicleCache) GetCount(ctx context.Context, keys ...any) (result int, err error) {
	cacheKey := v.buildKey(keys...)

	val, err := v.client.Get(ctx, cacheKey).Int()
	if err == redis.Nil {
		// Not found di cache bukan berarti error, return nil
		return 0, nil
	} else if err != nil {
		return 0, errors.LogAndReturnError(v.logger, err, "failed to get redis cache", "key", cacheKey)
	}

	return val, nil
}

// Set implements repository.VehicleCacheRepository.
func (v *vehicleCache) Set(ctx context.Context, ttl time.Duration, data any, keys ...any) error {
	cacheKey := v.buildKey(keys...)

	if err := v.client.Set(ctx, cacheKey, data, ttl).Err(); err != nil {
		return errors.LogAndReturnError(v.logger, err, "failed to set redis cache", "key", cacheKey)
	}
	return nil
}

// NewCustomerCache dengan prefix custom
func NewVehicleCache(logger *log.Logger, client *redis.Client) repository.VehicleCacheRepository {
	return &vehicleCache{
		client: client,
		logger: logger,
		prefix: "vehicles:",
	}
}
