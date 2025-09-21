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

type gormProviderVehicle struct {
	db     *gorm.DB
	logger *log.Logger
}

// CountVehicle implements repository.VehicleRepository.
func (g *gormProviderVehicle) CountVehicle(ctx context.Context, search string) (int, error) {
	var count int64
	query := g.db.WithContext(ctx).Model(&entities.Vehicle{})
	if search != "" {
		query = query.Where("brand LIKE ? OR model LIKE ? OR plate_no LIKE ? OR year LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	err := query.Count(&count).Error
	if err != nil {
		return 0, errors.ErrSomethingWrong(g.logger, err)
	}
	return int(count), nil
}

// SaveTx implements repository.VehicleRepository.
func (g *gormProviderVehicle) SaveBatchTx(ctx context.Context, tx any, m []entities.Vehicle) error {
	return errors.ErrDuplicateValue(g.logger, g.GetModelName(), tx.(*gorm.DB).WithContext(ctx).Create(&m).Error)
}

// Delete implements repository.VehicleRepository.
func (g *gormProviderVehicle) Delete(ctx context.Context, id uuid.UUID) error {
	return errors.ErrSomethingWrong(g.logger, g.db.WithContext(ctx).Where("id = ?", id).Updates(map[string]interface{}{"is_active": 0, "deleted_at": time.Now()}).Error)

}

// FindByCustomer implements repository.VehicleRepository.
func (g *gormProviderVehicle) FindByCustomer(ctx context.Context, customerID, limit, offset int, search string) (result base.SliceResult[*entities.Vehicle]) {
	query := g.db.WithContext(ctx).Preload("Devices").Preload("Sessions").Preload("Locations").Limit(limit).Offset(offset)
	if search != "" {
		query = query.Where("brand LIKE ? OR model LIKE ? OR plate_no LIKE ? OR year LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	if customerID != 0 {
		query = query.Where("customer_id = ?", customerID)
	}
	result.Error = errors.ErrRecordNotFound(g.logger, g.GetModelName(), query.Find(&result.Value).Error)
	return result
}

// FindByParam implements repository.VehicleRepository.
func (g *gormProviderVehicle) FindByParam(ctx context.Context, params map[string]interface{}) (result base.SingleResult[*entities.Vehicle]) {
	result.Error = errors.ErrRecordNotFound(g.logger, g.GetModelName(), g.db.WithContext(ctx).Preload("Devices").Preload("Sessions").Preload("Locations").Where(params).First(&result.Value).Error)
	return result
}

// GetConnection implements repository.VehicleRepository.
func (g *gormProviderVehicle) GetConnection() (T any) {
	return g.db
}

// GetModelName implements repository.VehicleRepository.
func (g *gormProviderVehicle) GetModelName() string {
	return "vehicles"
}

// Save implements repository.VehicleRepository.
func (g *gormProviderVehicle) Save(ctx context.Context, v *entities.Vehicle) error {
	return errors.ErrDuplicateValue(g.logger, g.GetModelName(), g.db.WithContext(ctx).Create(v).Error)
}

// Update implements repository.VehicleRepository.
func (g *gormProviderVehicle) Update(ctx context.Context, v *entities.Vehicle) error {
	return errors.ErrDuplicateValue(g.logger, g.GetModelName(), g.db.WithContext(ctx).Where("id = ?", v.ID).Updates(v).Error)
}

func NewGormVehicle(logger *log.Logger, db *gorm.DB) repository.VehicleRepository {
	return &gormProviderVehicle{db: db, logger: logger}
}
