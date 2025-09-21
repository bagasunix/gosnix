package postgres

import (
	"github.com/phuslu/log"
	"gorm.io/gorm"

	"github.com/bagasunix/gosnix/internal/domain/repository"
)

type Repositories interface {
	GetCustomer() repository.CustomerRepository
	GetHealth() repository.PostgresRepository
	GetVehicle() repository.VehicleRepository
	GetDeviceGPS() repository.DeviceRepository
	GetVehicleDevice() repository.VehicleDeviceRepository
	GetVehicleCategory() repository.VehicleCategoryRepository
}

type repo struct {
	customer      repository.CustomerRepository
	health        repository.PostgresRepository
	vehicle       repository.VehicleRepository
	device        repository.DeviceRepository
	vehicleDevice repository.VehicleDeviceRepository
	vehicleCat    repository.VehicleCategoryRepository
}

// GetVehicleCategory implements Repositories.
func (r *repo) GetVehicleCategory() repository.VehicleCategoryRepository {
	return r.vehicleCat
}

// GetVehicleDevice implements Repositories.
func (r *repo) GetVehicleDevice() repository.VehicleDeviceRepository {
	return r.vehicleDevice
}

// GetDeviceGPS implements Repositories.
func (r *repo) GetDeviceGPS() repository.DeviceRepository {
	return r.device
}

// GetVehicle implements Repositories.
func (r *repo) GetVehicle() repository.VehicleRepository {
	return r.vehicle
}

// GetCustomer implements Repositories.
func (r *repo) GetCustomer() repository.CustomerRepository {
	return r.customer
}

// GetHealth implements Repositories.
func (r *repo) GetHealth() repository.PostgresRepository {
	return r.health
}

func New(logger *log.Logger, db *gorm.DB) Repositories {
	return &repo{
		customer:      NewGormCustomer(logger, db),
		health:        NewHealthRepo(logger, db),
		vehicle:       NewGormVehicle(logger, db),
		device:        NewGormDevice(logger, db),
		vehicleDevice: NewGormVehicleDevice(logger, db),
		vehicleCat:    NewGormVehicleCategory(logger, db),
	}
}
