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

type gormProviderDevice struct {
	db     *gorm.DB
	logger *log.Logger
}

// CountCustomer implements repository.DeviceRepository.
func (g *gormProviderDevice) CountDevice(ctx context.Context, search string) (int, error) {
	var count int64
	query := g.db.WithContext(ctx).Model(&entities.DeviceGPS{})
	if search != "" {
		query = query.Where("imei LIKE ? OR brand LIKE ? OR model LIKE ? OR protocol LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	err := query.Count(&count).Error
	if err != nil {
		return 0, errors.ErrSomethingWrong(g.logger, err)
	}
	return int(count), nil
}

// Delete implements repository.DeviceRepository.
func (g *gormProviderDevice) Delete(ctx context.Context, id int) error {
	return errors.ErrSomethingWrong(g.logger, g.db.Model(&entities.DeviceGPS{}).WithContext(ctx).Where("id = ?", id).Updates(map[string]interface{}{"deleted_at": time.Now()}).Error)
}

// FindAll implements repository.DeviceRepository.
func (g *gormProviderDevice) FindAll(ctx context.Context, limit int, offset int, search string) (result base.SliceResult[*entities.DeviceGPS]) {
	query := g.db.WithContext(ctx).Preload("Vehicles").Limit(limit).Offset(offset)
	if search != "" {
		query = query.Where("imei LIKE ? OR brand LIKE ? OR model LIKE ? OR protocol LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	result.Error = errors.ErrRecordNotFound(g.logger, g.GetModelName(), query.Find(&result.Value).Error)
	return result
}

// FindByParams implements repository.DeviceRepository.
func (g *gormProviderDevice) FindByParam(ctx context.Context, params map[string]interface{}) (result base.SingleResult[*entities.DeviceGPS]) {
	result.Error = errors.ErrRecordNotFound(g.logger, g.GetModelName(), g.db.WithContext(ctx).Preload("Vehicles").Where(params).First(&result.Value).Error)
	return result
}

// GetConnection implements repository.DeviceRepository.
func (g *gormProviderDevice) GetConnection() (T any) {
	return g.db
}

// GetModelName implements repository.DeviceRepository.
func (g *gormProviderDevice) GetModelName() string {
	return "DeviceGPS"
}

// Save implements repository.DeviceRepository.
func (g *gormProviderDevice) Save(ctx context.Context, m *entities.DeviceGPS) error {
	return errors.ErrDuplicateValue(g.logger, g.GetModelName(), g.db.WithContext(ctx).Create(m).Error)
}

// SaveTx implements repository.DeviceRepository.
func (g *gormProviderDevice) SaveBatchTx(ctx context.Context, tx any, m []entities.DeviceGPS) error {
	return errors.ErrDuplicateValue(g.logger, g.GetModelName(), tx.(*gorm.DB).WithContext(ctx).Create(m).Error)
}

// Updates implements repository.DeviceRepository.
func (g *gormProviderDevice) Updates(ctx context.Context, id int, m *entities.DeviceGPS) error {
	return errors.ErrDuplicateValue(g.logger, g.GetModelName(), g.db.WithContext(ctx).Where("id = ?", id).Updates(m).Error)
}

func NewGormDevice(logger *log.Logger, db *gorm.DB) repository.DeviceRepository {
	return &gormProviderDevice{db: db, logger: logger}
}
