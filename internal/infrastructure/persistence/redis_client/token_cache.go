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

// tokenCache adalah implementasi caching token berbasis Redis
type tokenCache struct {
	client *redis.Client
	logger *log.Logger
	prefix string
}

// ---------------------------------------------------
// 🔹 Helper untuk membangun key
// ---------------------------------------------------
func (c *tokenCache) buildKey(parts ...any) string {
	key := c.prefix
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
func (t *tokenCache) Delete(ctx context.Context, keys ...any) error {
	cacheKey := t.buildKey(keys...)
	if err := t.client.Del(ctx, cacheKey).Err(); err != nil {
		return errors.LogAndReturnError(t.logger, err, "failed to delete redis cache", "key", cacheKey)
	}
	return nil
}

// ---------------------------------------------------
// 🔹 GET (string value)
// ---------------------------------------------------
func (t *tokenCache) Get(ctx context.Context, keys ...any) (result string, err error) {
	cacheKey := t.buildKey(keys...)

	val, err := t.client.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		return "", nil // tidak ditemukan
	} else if err != nil {
		return "", errors.LogAndReturnError(t.logger, err, "failed to get redis cache", "key", cacheKey)
	}
	return val, nil
}

// ---------------------------------------------------
// 🔹 SET (dengan TTL dan dukungan pipeline)
// ---------------------------------------------------
func (t *tokenCache) Set(ctx context.Context, ttl time.Duration, data any, keys ...any) error {
	cacheKey := t.buildKey(keys...)

	// Jika data bukan string, serialisasi ke JSON
	var val string
	switch d := data.(type) {
	case string:
		val = d
	default:
		b, err := json.Marshal(d)
		if err != nil {
			return errors.LogAndReturnError(t.logger, err, "failed to marshal cache data", "key", cacheKey)
		}
		val = string(b)
	}

	if err := t.client.Set(ctx, cacheKey, val, ttl).Err(); err != nil {
		return errors.LogAndReturnError(t.logger, err, "failed to set redis cache", "key", cacheKey)
	}
	return nil
}

// ---------------------------------------------------
// 🔹 GET Count (int value)
// ---------------------------------------------------
func (c *tokenCache) GetCount(ctx context.Context, keys ...any) (int, error) {
	cacheKey := c.buildKey(keys...)

	val, err := c.client.Get(ctx, cacheKey).Int()
	if err == redis.Nil {
		return 0, nil
	} else if err != nil {
		return 0, errors.LogAndReturnError(c.logger, err, "failed to get redis cache count", "key", cacheKey)
	}
	return val, nil
}

// ---------------------------------------------------
// 🔹 DELETE by pattern (gunakan SCAN + batching)
// ---------------------------------------------------
func (c *tokenCache) DeleteByPattern(ctx context.Context, pattern string) error {
	iter := c.client.Scan(ctx, 0, pattern, 100).Iterator() // batch 100 key per iterasi
	pipe := c.client.Pipeline()

	for iter.Next(ctx) {
		pipe.Del(ctx, iter.Val())
	}

	if _, err := pipe.Exec(ctx); err != nil {
		return errors.LogAndReturnError(c.logger, err, "failed to delete redis cache by pattern", "pattern", pattern)
	}

	if err := iter.Err(); err != nil {
		return errors.LogAndReturnError(c.logger, err, "failed to iterate scan", "pattern", pattern)
	}

	return nil
}

// NewTokenCache dengan prefix custom
func NewTokenCache(logger *log.Logger, client *redis.Client) repository.TokenCacheRepository {
	return &tokenCache{
		client: client,
		logger: logger,
		prefix: "tokens:",
	}
}
