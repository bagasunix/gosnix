package entities

import (
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type TrackingSession struct {
	ID            uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	VehicleID     uuid.UUID `gorm:"type:uuid;not null;index"`
	SessionName   string    `gorm:"size:255;not null"`
	StartTime     time.Time `gorm:"not null"`
	EndTime       *time.Time
	TotalDistance float64        `gorm:"type:decimal(10,3);default:0.0"`
	TotalDuration int            `gorm:"default:0"`
	CreatedBy     int            `gorm:"not null"`
	CreatedAt     time.Time      `gorm:"autoCreateTime"`
	UpdatedAt     *time.Time     `gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`

	Vehicle       Vehicle          `gorm:"foreignKey:VehicleID;references:ID"`
	CreatedByUser User             `gorm:"foreignKey:CreatedBy;references:ID"`
	Locations     []LocationUpdate `gorm:"foreignKey:SessionID;references:ID"`
}

func (TrackingSession) TableName() string {
	return "tracking_sessions"
}
