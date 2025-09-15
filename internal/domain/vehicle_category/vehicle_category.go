package vehicle_category

import (
	"time"

	"github.com/gofrs/uuid"
)

type VehicleCategory struct {
	ID          uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name        string     `gorm:"size:100;not null"`
	Description string     `gorm:"type:text"`
	CreatedAt   time.Time  `gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime"`
	DeletedAt   *time.Time `gorm:"index"`
}

func (VehicleCategory) TableName() string {
	return "vehicle_categories"
}
