package locationupdate

import (
	"time"

	"github.com/gofrs/uuid"
)

type LocationUpdate struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	VehicleID  uuid.UUID `gorm:"type:uuid;not null"`
	SessionID  uuid.UUID `gorm:"type:uuid;not null"`
	Latitude   float64   `gorm:"type:decimal(10,8);not null"`
	Longitude  float64   `gorm:"type:decimal(11,8);not null"`
	ReceivedAt time.Time `gorm:"autoCreateTime"`
}

func (LocationUpdate) TableName() string {
	return "location_updates"
}
