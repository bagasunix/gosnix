package repository

import (
	"context"

	"github.com/bagasunix/gosnix/internal/domain/base"
	"github.com/bagasunix/gosnix/internal/domain/entities"
)

type DeviceCommon interface {
	Save(ctx context.Context, m *entities.DeviceGPS) error
	SaveBatchTx(ctx context.Context, tx any, m []entities.DeviceGPS) error
	Delete(ctx context.Context, id int) error
	Updates(ctx context.Context, id int, m *entities.DeviceGPS) error
	FindAll(ctx context.Context, limit int, offset int, search string) (result base.SliceResult[*entities.DeviceGPS])
	CountCustomer(ctx context.Context, search string) (int, error)

	FindByParam(ctx context.Context, params map[string]interface{}) (result base.SingleResult[*entities.DeviceGPS])
}
type DeviceRepository interface {
	base.Repository
	DeviceCommon
}
