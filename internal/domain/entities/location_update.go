package entities

import (
	"time"

	"github.com/gofrs/uuid"
)

type LocationUpdate struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	VehicleID  uuid.UUID `gorm:"type:uuid;not null;index"`
	SessionID  uuid.UUID `gorm:"type:uuid;not null;index"`
	Latitude   float64   `gorm:"type:decimal(10,8);not null"`
	Longitude  float64   `gorm:"type:decimal(11,8);not null"`
	Speed      float64   `gorm:"type:decimal(6,2)"`
	Heading    int       `gorm:"type:smallint"`
	ReceivedAt time.Time `gorm:"autoCreateTime"`

	Vehicle Vehicle         `gorm:"foreignKey:VehicleID;references:ID"`
	Session TrackingSession `gorm:"foreignKey:SessionID;references:ID"`
}

func (LocationUpdate) TableName() string {
	return "location_updates"
}
