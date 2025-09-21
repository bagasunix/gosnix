package postgres

import (
	"context"
	"time"

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

// Delete implements repository.VehicleDeviceRepository.
func (g *gormProviderDeviceVehicle) Delete(ctx context.Context, id int) error {
	return errors.ErrSomethingWrong(g.logger, g.db.Model(&entities.VehicleDevice{}).WithContext(ctx).Where("id = ?", id).Updates(map[string]any{"is_active": 0, "deleted_at": time.Now(), "end_time": time.Now()}).Error)
}

// FindByParams implements repository.VehicleDeviceRepository.
func (g *gormProviderDeviceVehicle) FindByParams(ctx context.Context, params map[string]interface{}) (result base.SingleResult[*entities.VehicleDevice]) {
	result.Error = errors.ErrRecordNotFound(g.logger, g.GetModelName(), g.db.WithContext(ctx).Where(params).First(&result.Value).Error)
	return result
}

// GetConnection implements repository.VehicleDeviceRepository.
func (g *gormProviderDeviceVehicle) GetConnection() (T any) {
	return g.db
}

// GetModelName implements repository.VehicleDeviceRepository.
func (g *gormProviderDeviceVehicle) GetModelName() string {
	return "VehicleDevice"
}

// Save implements repository.VehicleDeviceRepository.
func (g *gormProviderDeviceVehicle) Save(ctx context.Context, m *entities.VehicleDevice) error {
	return errors.ErrDuplicateValue(g.logger, g.GetModelName(), g.db.WithContext(ctx).Create(m).Error)
}

// SaveTx implements repository.VehicleDeviceRepository.
func (g *gormProviderDeviceVehicle) SaveBatchTx(ctx context.Context, tx any, m []entities.VehicleDevice) error {
	return errors.ErrDuplicateValue(g.logger, g.GetModelName(), tx.(*gorm.DB).WithContext(ctx).Create(m).Error)
}

// Updates implements repository.VehicleDeviceRepository.
func (g *gormProviderDeviceVehicle) Updates(ctx context.Context, id int, m *entities.VehicleDevice) error {
	return errors.ErrDuplicateValue(g.logger, g.GetModelName(), g.db.WithContext(ctx).Where("id = ?", id).Updates(m).Error)
}

func NewGormVehicleDevice(logger *log.Logger, db *gorm.DB) repository.VehicleDeviceRepository {
	return &gormProviderDeviceVehicle{db: db, logger: logger}
}
