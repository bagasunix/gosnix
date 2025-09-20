package entities

import (
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type VehicleCategory struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name        string         `gorm:"size:100;uniqueIndex;not null"`
	Description string         `gorm:"type:text"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   *time.Time     `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`

	Vehicles []Vehicle `gorm:"foreignKey:CategoryID;references:ID"`
}

func (VehicleCategory) TableName() string {
	return "vehicle_categories"
}
