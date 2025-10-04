package redis_client

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/phuslu/log"
	"github.com/redis/go-redis/v9"

	"github.com/bagasunix/gosnix/internal/domain/repository"
	"github.com/bagasunix/gosnix/pkg/errors"
)

// vehicleCache adalah implementasi caching vehicle berbasis Redis
type vehicleCache struct {
	client *redis.Client
	logger *log.Logger
	prefix string
}

// ---------------------------------------------------
// 🔹 Helper untuk membangun key
// ---------------------------------------------------
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

// ---------------------------------------------------
// 🔹 DELETE satu key
// ---------------------------------------------------
func (v *vehicleCache) Delete(ctx context.Context, keys ...any) error {
	cacheKey := v.buildKey(keys...)
	if err := v.client.Del(ctx, cacheKey).Err(); err != nil {
		return errors.LogAndReturnError(v.logger, err, "failed to delete redis cache", "key", cacheKey)
	}
	return nil
}

// ---------------------------------------------------
// 🔹 DELETE by pattern (gunakan SCAN + batching)
// ---------------------------------------------------
func (v *vehicleCache) DeleteByPattern(ctx context.Context, pattern string) error {
	iter := v.client.Scan(ctx, 0, pattern, 100).Iterator() // batch 100 key per iterasi
	pipe := v.client.Pipeline()

	for iter.Next(ctx) {
		pipe.Del(ctx, iter.Val())
	}

	if _, err := pipe.Exec(ctx); err != nil {
		return errors.LogAndReturnError(v.logger, err, "failed to delete redis cache by pattern", "pattern", pattern)
	}

	if err := iter.Err(); err != nil {
		return errors.LogAndReturnError(v.logger, err, "failed to iterate scan", "pattern", pattern)
	}

	return nil
}

// ---------------------------------------------------
// 🔹 GET (string value)
// ---------------------------------------------------
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

// ---------------------------------------------------
// 🔹 GET Count (int value)
// ---------------------------------------------------
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

// ---------------------------------------------------
// 🔹 SET (dengan TTL dan dukungan pipeline)
// ---------------------------------------------------
func (v *vehicleCache) Set(ctx context.Context, ttl time.Duration, data any, keys ...any) error {
	cacheKey := v.buildKey(keys...)

	// Jika data bukan string, serialisasi ke JSON
	var val string
	switch d := data.(type) {
	case string:
		val = d
	default:
		b, err := json.Marshal(d)
		if err != nil {
			return errors.LogAndReturnError(v.logger, err, "failed to marshal cache data", "key", cacheKey)
		}
		val = string(b)
	}

	if err := v.client.Set(ctx, cacheKey, val, ttl).Err(); err != nil {
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
