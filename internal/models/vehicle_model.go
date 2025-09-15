package models

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/bagasunix/gosnix/internal/domain/vehicle"
)

type VehicleModel struct {
	ID             uuid.UUID
	CategoryID     uuid.UUID
	PlateNo        string
	Model          string
	Years          string
	Color          string
	VIN            string
	DeviceID       string
	DeviceType     string
	FuelType       string
	EngineCapacity float64
	MaxSpeed       int
	IsActive       int
	CreatedBy      int
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time

	Category VehicleCategoryModel `gorm:"foreignKey:CategoryID"`
	// Merchant          Merchant          `gorm:"foreignKey:MerchantID"`
	CreatedByCustomer CustomerModel           `gorm:"foreignKey:CreatedBy"`
	Sessions          []TrackingSessionMModel `gorm:"foreignKey:VehicleID"`
	Locations         []LocationUpdateModels  `gorm:"foreignKey:VehicleID"`
}

func (v *VehicleModel) ToDomain() *vehicle.Vehicle {
	return &vehicle.Vehicle{
		ID:             v.ID,
		CategoryID:     v.CategoryID,
		PlateNo:        v.PlateNo,
		Model:          v.Model,
		Years:          v.Years,
		Color:          v.Color,
		VIN:            v.VIN,
		DeviceID:       v.DeviceID,
		DeviceType:     v.DeviceType,
		FuelType:       v.FuelType,
		EngineCapacity: v.EngineCapacity,
		MaxSpeed:       v.MaxSpeed,
		IsActive:       v.IsActive,
		CreatedBy:      v.CreatedBy,
		CreatedAt:      v.CreatedAt,
		UpdatedAt:      v.UpdatedAt,
		DeletedAt:      v.DeletedAt,
	}

}
