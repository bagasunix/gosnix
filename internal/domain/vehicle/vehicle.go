package vehicle

import (
	"time"

	"github.com/gofrs/uuid"
)

type Vehicle struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CategoryID uuid.UUID `gorm:"type:uuid;not null"`
	// MerchantID     uuid.UUID `gorm:"type:uuid;not null"`
	PlateNo        string  `gorm:"uniqueIndex;size:20"`
	Model          string  `gorm:"size:100"`
	Years          string  `gorm:"size:10"`
	Color          string  `gorm:"size:50"`
	VIN            string  `gorm:"size:50"`
	DeviceID       string  `gorm:"uniqueIndex;size:100"`
	DeviceType     string  `gorm:"size:50"`
	FuelType       string  `gorm:"size:50"`
	EngineCapacity float64 `gorm:"type:decimal(6,2)"`
	MaxSpeed       int
	IsActive       int
	CreatedBy      int
	CreatedAt      time.Time  `gorm:"autoCreateTime"`
	UpdatedAt      time.Time  `gorm:"autoUpdateTime"`
	DeletedAt      *time.Time `gorm:"index"`
}

func (Vehicle) TableName() string {
	return "vehicles"
}
