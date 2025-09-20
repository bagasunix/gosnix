package entities

import (
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type DeviceGPS struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	IMEI      string         `gorm:"size:15;uniqueIndex;not null"`
	Brand     string         `gorm:"size:100;not null"`
	Model     string         `gorm:"size:100;not null"`
	Protocol  string         `gorm:"size:20;not null"` // TCP, UDP, HTTP, MQTT
	SecretKey string         `gorm:"size:255;not null"`
	CreatedBy int            `gorm:"not null"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Vehicles      []VehicleDevice `gorm:"foreignKey:DeviceID;references:ID"`
	CreatedByUser User            `gorm:"foreignKey:CreatedBy;references:ID"`
}

func (DeviceGPS) TableName() string {
	return "device_gps"
}
