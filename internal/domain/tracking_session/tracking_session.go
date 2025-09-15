package trackingsession

import (
	"time"

	"github.com/gofrs/uuid"
)

type TrackingSession struct {
	ID            uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	VehicleID     uuid.UUID `gorm:"type:uuid;not null"`
	SessionName   string    `gorm:"size:255"`
	StartTime     time.Time `gorm:"autoCreateTime"`
	EndTime       *time.Time
	Status        string  `gorm:"size:50;default:ACTIVE"`
	TotalDistance float64 `gorm:"type:decimal(10,3)"`
	TotalDuration int
	CreatedBy     int
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
}

func (TrackingSession) TableName() string {
	return "tracking_sessions"
}
