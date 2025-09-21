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

type gormProviderDevice struct {
	db     *gorm.DB
	logger *log.Logger
}

// CountCustomer implements repository.DeviceRepository.
func (g *gormProviderDevice) CountCustomer(ctx context.Context, search string) (int, error) {
	panic("unimplemented")
}

// Delete implements repository.DeviceRepository.
func (g *gormProviderDevice) Delete(ctx context.Context, id int) error {
	panic("unimplemented")
}

// FindAll implements repository.DeviceRepository.
func (g *gormProviderDevice) FindAll(ctx context.Context, limit int, offset int, search string) (result base.SliceResult[*entities.DeviceGPS]) {
	panic("unimplemented")
}

// FindByParams implements repository.DeviceRepository.
func (g *gormProviderDevice) FindByParam(ctx context.Context, params map[string]interface{}) (result base.SingleResult[*entities.DeviceGPS]) {
	panic("unimplemented")
}

// GetConnection implements repository.DeviceRepository.
func (g *gormProviderDevice) GetConnection() (T any) {
	panic("unimplemented")
}

// GetModelName implements repository.DeviceRepository.
func (g *gormProviderDevice) GetModelName() string {
	panic("unimplemented")
}

// Save implements repository.DeviceRepository.
func (g *gormProviderDevice) Save(ctx context.Context, m *entities.DeviceGPS) error {
	panic("unimplemented")
}

// SaveTx implements repository.DeviceRepository.
func (g *gormProviderDevice) SaveBatchTx(ctx context.Context, tx any, m []entities.DeviceGPS) error {
	return errors.ErrDuplicateValue(g.logger, g.GetModelName(), tx.(*gorm.DB).WithContext(ctx).Create(m).Error)
}

// Updates implements repository.DeviceRepository.
func (g *gormProviderDevice) Updates(ctx context.Context, id int, m *entities.DeviceGPS) error {
	panic("unimplemented")
}

func NewGormDevice(logger *log.Logger, db *gorm.DB) repository.DeviceRepository {
	return &gormProviderDevice{db: db, logger: logger}
}
