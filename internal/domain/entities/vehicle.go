package entities

import (
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type Vehicle struct {
	ID              uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CategoryID      uuid.UUID      `gorm:"type:uuid;not null"`
	CustomerID      int            `gorm:"not null"`
	PlateNo         string         `gorm:"uniqueIndex;size:20;not null"`
	Model           string         `gorm:"size:100;not null"`
	Brand           string         `gorm:"size:100;not null"`
	ManufactureYear int            `gorm:"size:4;not null"`
	Color           string         `gorm:"size:50;not null"`
	VIN             string         `gorm:"size:50;uniqueIndex;not null"`
	FuelType        string         `gorm:"size:50;not null"`
	EngineCapacity  float64        `gorm:"type:decimal(6,2);not null"`
	MaxSpeed        int            `gorm:"not null"`
	IsActive        int            `gorm:"default:0;not null"`
	CreatedBy       int            `gorm:"not null"`
	CreatedAt       time.Time      `gorm:"autoCreateTime"`
	UpdatedAt       *time.Time     `gorm:"autoUpdateTime"`
	DeletedAt       gorm.DeletedAt `gorm:"index"`

	Category  VehicleCategory   `gorm:"foreignKey:CategoryID;references:ID"`
	Customer  Customer          `gorm:"foreignKey:CustomerID;references:ID"`
	Devices   []VehicleDevice   `gorm:"foreignKey:VehicleID;references:ID"`
	Sessions  []TrackingSession `gorm:"foreignKey:VehicleID;references:ID"`
	Locations []LocationUpdate  `gorm:"foreignKey:VehicleID;references:ID"`
}

func (Vehicle) TableName() string {
	return "vehicles"
}
