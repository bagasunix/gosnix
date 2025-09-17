package entities

import (
	"time"

	"github.com/gofrs/uuid"
)

type Vehicle struct {
	ID             uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CategoryID     uuid.UUID `gorm:"type:uuid;not null"`
	CustomerID     int       `gorm:"not null"`
	PlateNo        string    `gorm:"uniqueIndex;size:20"`
	Model          string    `gorm:"size:100"`
	Brand          string    `gorm:"size:100"`
	Years          string    `gorm:"size:10"`
	Color          string    `gorm:"size:50"`
	VIN            string    `gorm:"size:50"`
	DeviceID       string    `gorm:"uniqueIndex;size:100"`
	DeviceType     string    `gorm:"size:50"`
	FuelType       string    `gorm:"size:50"`
	EngineCapacity float64   `gorm:"type:decimal(6,2)"`
	MaxSpeed       int
	Year           string
	IsActive       int
	CreatedBy      int
	CreatedAt      time.Time  `gorm:"autoCreateTime"`
	UpdatedAt      time.Time  `gorm:"autoUpdateTime"`
	DeletedAt      *time.Time `gorm:"index"`

	Category  VehicleCategory   `gorm:"foreignKey:CategoryID"`
	Sessions  []TrackingSession `gorm:"foreignKey:VehicleID"`
	Locations []LocationUpdate  `gorm:"foreignKey:VehicleID"`
}

func (Vehicle) TableName() string {
	return "vehicles"
}
