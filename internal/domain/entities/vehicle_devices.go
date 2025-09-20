package entities

import (
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type VehicleDevice struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	VehicleID uuid.UUID `gorm:"type:uuid;not null;index"`
	DeviceID  uuid.UUID `gorm:"type:uuid;not null;index"`
	StartTime time.Time `gorm:"not null"`
	EndTime   *time.Time
	IsActive  int            `gorm:"default:0;not null"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt *time.Time     `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Vehicle Vehicle   `gorm:"foreignKey:VehicleID;references:ID"`
	Device  DeviceGPS `gorm:"foreignKey:DeviceID;references:ID"`
}

func (VehicleDevice) TableName() string {
	return "vehicle_devices"
}
