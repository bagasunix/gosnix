package repository

import (
	"context"

	"github.com/bagasunix/gosnix/internal/domain/base"
	"github.com/bagasunix/gosnix/internal/domain/entities"
)

type VehicleDeviceCommon interface {
	Save(ctx context.Context, m *entities.VehicleDevice) error
	SaveBatchTx(ctx context.Context, tx any, m []entities.VehicleDevice) error
	Delete(ctx context.Context, id int) error
	Updates(ctx context.Context, id int, m *entities.VehicleDevice) error
	FindAll(ctx context.Context, limit int, offset int, search string) (result base.SliceResult[*entities.VehicleDevice])
	CountCustomer(ctx context.Context, search string) (int, error)

	FindByParams(ctx context.Context, params map[string]interface{}) (result base.SingleResult[*entities.VehicleDevice])
}
type VehicleDeviceRepository interface {
	base.Repository
	VehicleDeviceCommon
}
