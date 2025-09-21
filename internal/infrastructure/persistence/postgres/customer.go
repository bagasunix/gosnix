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

type gormProviderCustomer struct {
	db     *gorm.DB
	logger *log.Logger
}

// FindByPhoneOrEmail implements repository.CustomerRepository.
func (g *gormProviderCustomer) FindByPhoneOrEmail(ctx context.Context, phone string, email string) (result base.SingleResult[*entities.Customer]) {
	result.Error = errors.ErrRecordNotFound(g.logger, g.GetModelName(), g.db.WithContext(ctx).Preload("Vehicles").Where("phone = ? OR email = ?", phone, email).First(&result.Value).Error)
	return result
}

// CountCustomer implements entities.Customer .Repository.
func (g *gormProviderCustomer) CountCustomer(ctx context.Context, search string) (int, error) {
	var count int64
	query := g.db.WithContext(ctx).Model(&entities.Customer{})
	if search != "" {
		query = query.Where("name LIKE ? OR email LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	err := query.Count(&count).Error
	if err != nil {
		return 0, errors.ErrSomethingWrong(g.logger, err)
	}
	return int(count), nil
}

// Delete implements entities.Repository.
func (g *gormProviderCustomer) Delete(ctx context.Context, id int) error {
	return errors.ErrSomethingWrong(g.logger, g.db.Model(&entities.Customer{}).WithContext(ctx).Where("id = ?", id).Updates(map[string]interface{}{"is_active": 0, "deleted_at": time.Now()}).Error)
}

// FindAll implements entities.Repository.
func (g *gormProviderCustomer) FindAll(ctx context.Context, limit int, offset int, search string) (result base.SliceResult[*entities.Customer]) {
	query := g.db.WithContext(ctx).Preload("Vehicles.Devices.Device").Limit(limit).Offset(offset)
	if search != "" {
		query = query.Where("name LIKE ? OR email LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	result.Error = errors.ErrRecordNotFound(g.logger, g.GetModelName(), query.Find(&result.Value).Error)
	return result
}

// FindByParams implements entities.Repository.
func (g *gormProviderCustomer) FindByParam(ctx context.Context, params map[string]interface{}) (result base.SingleResult[*entities.Customer]) {
	result.Error = errors.ErrRecordNotFound(g.logger, g.GetModelName(), g.db.WithContext(ctx).Preload("Vehicles.Devices.Device").Where(params).First(&result.Value).Error)
	return result
}

func (g *gormProviderCustomer) GetConnection() (T any) {
	return g.db
}
func (g *gormProviderCustomer) GetModelName() string {
	return "customers"
}

// Save implements entities.Repository.
func (g *gormProviderCustomer) Save(ctx context.Context, m *entities.Customer) error {
	return errors.ErrDuplicateValue(g.logger, g.GetModelName(), g.db.WithContext(ctx).Create(m).Error)
}

// SaveTx implements entities.Repository.
func (g *gormProviderCustomer) SaveTx(ctx context.Context, tx any, m *entities.Customer) error {
	return errors.ErrDuplicateValue(g.logger, g.GetModelName(), tx.(*gorm.DB).WithContext(ctx).Create(m).Error)
}

// Updates implements entities.Repository.
func (g *gormProviderCustomer) Updates(ctx context.Context, id int, m *entities.Customer) error {
	return errors.ErrDuplicateValue(g.logger, g.GetModelName(), g.db.WithContext(ctx).Where("id = ?", id).Updates(m).Error)
}

func NewGormCustomer(logger *log.Logger, db *gorm.DB) repository.CustomerRepository {
	return &gormProviderCustomer{db: db, logger: logger}
}
