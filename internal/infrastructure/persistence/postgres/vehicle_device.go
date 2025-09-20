package postgres

import (
	"context"

	"github.com/phuslu/log"
	"gorm.io/gorm"

	"github.com/bagasunix/gosnix/internal/domain/base"
	"github.com/bagasunix/gosnix/internal/domain/entities"
	"github.com/bagasunix/gosnix/internal/domain/repository"
	"github.com/bagasunix/gosnix/pkg/errors"
)

type gormProviderDeviceVehicle struct {
	db     *gorm.DB
	logger *log.Logger
}

// CountCustomer implements repository.VehicleDeviceRepository.
func (g *gormProviderDeviceVehicle) CountCustomer(ctx context.Context, search string) (int, error) {
	panic("unimplemented")
}

// Delete implements repository.VehicleDeviceRepository.
func (g *gormProviderDeviceVehicle) Delete(ctx context.Context, id int) error {
	panic("unimplemented")
}

// FindAll implements repository.VehicleDeviceRepository.
func (g *gormProviderDeviceVehicle) FindAll(ctx context.Context, limit int, offset int, search string) (result base.SliceResult[*entities.VehicleDevice]) {
	panic("unimplemented")
}

// FindByParams implements repository.VehicleDeviceRepository.
func (g *gormProviderDeviceVehicle) FindByParams(ctx context.Context, params map[string]interface{}) (result base.SingleResult[*entities.VehicleDevice]) {
	panic("unimplemented")
}

// GetConnection implements repository.VehicleDeviceRepository.
func (g *gormProviderDeviceVehicle) GetConnection() (T any) {
	panic("unimplemented")
}

// GetModelName implements repository.VehicleDeviceRepository.
func (g *gormProviderDeviceVehicle) GetModelName() string {
	panic("unimplemented")
}

// Save implements repository.VehicleDeviceRepository.
func (g *gormProviderDeviceVehicle) Save(ctx context.Context, m *entities.VehicleDevice) error {
	panic("unimplemented")
}

// SaveTx implements repository.VehicleDeviceRepository.
func (g *gormProviderDeviceVehicle) SaveBatchTx(ctx context.Context, tx any, m []entities.VehicleDevice) error {
	return errors.ErrDuplicateValue(g.logger, g.GetModelName(), tx.(*gorm.DB).WithContext(ctx).Create(m).Error)
}

// Updates implements repository.VehicleDeviceRepository.
func (g *gormProviderDeviceVehicle) Updates(ctx context.Context, id int, m *entities.VehicleDevice) error {
	panic("unimplemented")
}

func NewGormVehicleDevice(logger *log.Logger, db *gorm.DB) repository.VehicleDeviceRepository {
	return &gormProviderDeviceVehicle{db: db, logger: logger}
}
