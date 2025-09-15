package models

import (
	"time"

	"github.com/gofrs/uuid"

	locationupdate "github.com/bagasunix/gosnix/internal/domain/location_update"
)

type LocationUpdateModels struct {
	ID         uuid.UUID
	VehicleID  uuid.UUID
	SessionID  uuid.UUID
	Latitude   float64
	Longitude  float64
	ReceivedAt time.Time

	Vehicle VehicleModel          `gorm:"foreignKey:VehicleID"`
	Session TrackingSessionMModel `gorm:"foreignKey:SessionID"`
}

func (lu *LocationUpdateModels) ToDomain() *locationupdate.LocationUpdate {
	return &locationupdate.LocationUpdate{
		ID:         lu.ID,
		VehicleID:  lu.VehicleID,
		SessionID:  lu.SessionID,
		Latitude:   lu.Latitude,
		Longitude:  lu.Longitude,
		ReceivedAt: lu.ReceivedAt,
	}
}
