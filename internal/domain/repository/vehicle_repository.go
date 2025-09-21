package repository

import (
	"context"
	"time"

	"github.com/gofrs/uuid"

	"github.com/bagasunix/gosnix/internal/domain/base"
	"github.com/bagasunix/gosnix/internal/domain/entities"
)

type VehicleCommon interface {
	Save(ctx context.Context, v *entities.Vehicle) error
	SaveBatchTx(ctx context.Context, tx any, m []entities.Vehicle) error
	FindByCustomer(ctx context.Context, customerID, limit, offset int, search string) (result base.SliceResult[*entities.Vehicle])
	FindByParam(ctx context.Context, params map[string]any) (result base.SingleResult[*entities.Vehicle])
	Update(ctx context.Context, v *entities.Vehicle) error
	Delete(ctx context.Context, id uuid.UUID) error
	CountVehicle(ctx context.Context, search string) (int, error)
}
type VehicleRepository interface {
	base.Repository
	VehicleCommon
}

type VehicleCacheRepository interface {
	Set(ctx context.Context, ttl time.Duration, data any, keys ...any) error
	Get(ctx context.Context, keys ...any) (result *string, err error)
	GetCount(ctx context.Context, keys ...any) (result int, err error)
	Delete(ctx context.Context, keys ...any) error
	DeleteByPattern(ctx context.Context, pattern string) error
}
