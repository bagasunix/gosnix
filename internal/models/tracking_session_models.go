package models

import (
	"time"

	"github.com/gofrs/uuid"

	trackingsession "github.com/bagasunix/gosnix/internal/domain/tracking_session"
)

type TrackingSessionMModel struct {
	ID            uuid.UUID
	VehicleID     uuid.UUID
	SessionName   string
	StartTime     time.Time
	EndTime       *time.Time
	Status        string
	TotalDistance float64
	TotalDuration int
	CreatedBy     int
	CreatedAt     time.Time
	UpdatedAt     time.Time

	Vehicle       VehicleModel           `gorm:"foreignKey:VehicleID"`
	CreatedByUser CustomerModel          `gorm:"foreignKey:CreatedBy"`
	Locations     []LocationUpdateModels `gorm:"foreignKey:SessionID"`
}

func (ts *TrackingSessionMModel) ToDomain() *trackingsession.TrackingSession {
	return &trackingsession.TrackingSession{
		ID:            ts.Vehicle.ID,
		VehicleID:     ts.VehicleID,
		SessionName:   ts.SessionName,
		StartTime:     ts.StartTime,
		EndTime:       ts.EndTime,
		Status:        ts.Status,
		TotalDistance: ts.TotalDistance,
		TotalDuration: ts.TotalDuration,
		CreatedBy:     ts.CreatedBy,
		CreatedAt:     ts.CreatedAt,
		UpdatedAt:     ts.UpdatedAt,
	}
}
