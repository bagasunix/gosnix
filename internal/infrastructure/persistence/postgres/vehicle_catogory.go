package postgres

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
	"github.com/phuslu/log"
	"gorm.io/gorm"

	"github.com/bagasunix/gosnix/internal/domain/base"
	"github.com/bagasunix/gosnix/internal/domain/entities"
	"github.com/bagasunix/gosnix/internal/domain/repository"
	"github.com/bagasunix/gosnix/pkg/errors"
)

type gormProviderVehicleCategory struct {
	db     *gorm.DB
	logger *log.Logger
}

// CountVehicleCategory implements repository.VehicleCategoryRepository.
func (g *gormProviderVehicleCategory) CountVehicleCategory(ctx context.Context, search string) (int, error) {
	var count int64
	query := g.db.WithContext(ctx).Model(&entities.VehicleCategory{})
	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}
	err := query.Count(&count).Error
	if err != nil {
		return 0, errors.ErrSomethingWrong(g.logger, err)
	}
	return int(count), nil
}

// Delete implements repository.VehicleCategoryRepository.
func (g *gormProviderVehicleCategory) Delete(ctx context.Context, id uuid.UUID) error {
	return errors.ErrSomethingWrong(g.logger, g.db.Model(&entities.VehicleCategory{}).WithContext(ctx).Where("id = ?", id).Updates(map[string]interface{}{"deleted_at": time.Now()}).Error)
}

// FindAll implements repository.VehicleCategoryRepository.
func (g *gormProviderVehicleCategory) FindAll(ctx context.Context, limit int, offset int, search string) (result base.SliceResult[*entities.VehicleCategory]) {
	query := g.db.WithContext(ctx).Preload("Vehicles").Limit(limit).Offset(offset)
	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}
	result.Error = errors.ErrRecordNotFound(g.logger, g.GetModelName(), query.Find(&result.Value).Error)
	return result
}

// FindByParam implements repository.VehicleCategoryRepository.
func (g *gormProviderVehicleCategory) FindByParam(ctx context.Context, param map[string]any) (result base.SingleResult[*entities.VehicleCategory]) {
	result.Error = errors.ErrRecordNotFound(g.logger, g.GetModelName(), g.db.WithContext(ctx).Preload("Vehicles").Where(param).First(&result.Value).Error)
	return result
}

// GetConnection implements repository.VehicleCategoryRepository.
func (g *gormProviderVehicleCategory) GetConnection() (T any) {
	return g.db
}

// GetModelName implements repository.VehicleCategoryRepository.
func (g *gormProviderVehicleCategory) GetModelName() string {
	return "VehicleCategory"
}

// Save implements repository.VehicleCategoryRepository.
func (g *gormProviderVehicleCategory) Save(ctx context.Context, vc *entities.VehicleCategory) error {
	return errors.ErrDuplicateValue(g.logger, g.GetModelName(), g.db.WithContext(ctx).Create(vc).Error)
}

// Update implements repository.VehicleCategoryRepository.
func (g *gormProviderVehicleCategory) Update(ctx context.Context, vc *entities.VehicleCategory) error {
	return errors.ErrDuplicateValue(g.logger, g.GetModelName(), g.db.WithContext(ctx).Where("id = ?", vc.ID).Updates(vc).Error)
}

func NewGormVehicleCategory(logger *log.Logger, db *gorm.DB) repository.VehicleCategoryRepository {
	return &gormProviderVehicleCategory{db: db, logger: logger}
}
