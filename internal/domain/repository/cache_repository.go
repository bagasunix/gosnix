package repository

import (
	"context"
	"time"

	"github.com/bagasunix/gosnix/internal/domain/entities"
)

// Repository Redis
type CustomerCacheRepository interface {
	Set(ctx context.Context, ttl time.Duration, data any, keys ...any) error
	GetWithValue(ctx context.Context, keys ...any) (result *entities.Customer, err error)
	Get(ctx context.Context, keys ...any) (result string, err error)
	GetCount(ctx context.Context, keys ...any) (result int, err error)
	Delete(ctx context.Context, keys ...any) error
	DeleteByPattern(ctx context.Context, pattern string) error
}

type VehicleCacheRepository interface {
	Set(ctx context.Context, ttl time.Duration, data any, keys ...any) error
	Get(ctx context.Context, keys ...any) (result *string, err error)
	GetCount(ctx context.Context, keys ...any) (result int, err error)
	Delete(ctx context.Context, keys ...any) error
	DeleteByPattern(ctx context.Context, pattern string) error
}

type TokenCacheRepository interface {
	Set(ctx context.Context, ttl time.Duration, data any, keys ...any) error
	Get(ctx context.Context, keys ...any) (result string, err error)
	Delete(ctx context.Context, keys ...any) error
	GetCount(ctx context.Context, keys ...any) (int, error)
	DeleteByPattern(ctx context.Context, pattern string) error
}
