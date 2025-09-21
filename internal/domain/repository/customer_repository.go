package repository

import (
	"context"
	"time"

	"github.com/bagasunix/gosnix/internal/domain/base"
	"github.com/bagasunix/gosnix/internal/domain/entities"
)

// Repository Database
type CustomerCommon interface {
	Save(ctx context.Context, m *entities.Customer) error
	SaveTx(ctx context.Context, tx any, m *entities.Customer) error
	Delete(ctx context.Context, id int) error
	Updates(ctx context.Context, id int, m *entities.Customer) error
	FindAll(ctx context.Context, limit int, offset int, search string) (result base.SliceResult[*entities.Customer])
	CountCustomer(ctx context.Context, search string) (int, error)

	FindByParam(ctx context.Context, params map[string]interface{}) (result base.SingleResult[*entities.Customer])
	FindByPhoneOrEmail(ctx context.Context, phone, email string) (result base.SingleResult[*entities.Customer])
}
type CustomerRepository interface {
	base.Repository
	CustomerCommon
}

// Repository Redis
type CustomerCacheRepository interface {
	Set(ctx context.Context, ttl time.Duration, data any, keys ...any) error
	GetWithValue(ctx context.Context, keys ...any) (result *entities.Customer, err error)
	Get(ctx context.Context, keys ...any) (result *string, err error)
	GetCount(ctx context.Context, keys ...any) (result int, err error)
	Delete(ctx context.Context, keys ...any) error
	DeleteByPattern(ctx context.Context, pattern string) error
}
