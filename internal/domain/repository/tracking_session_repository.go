package repository

import (
	"context"

	"github.com/gofrs/uuid"

	"github.com/bagasunix/gosnix/internal/domain/base"
	"github.com/bagasunix/gosnix/internal/domain/entities"
)

type TrackingSessionCommon interface {
	Save(ctx context.Context, vc *entities.TrackingSession) error
	FindByParam(ctx context.Context, param map[string]any) (result base.SingleResult[*entities.TrackingSession])
	FindByVehicle(ctx context.Context, vehicleID uuid.UUID) (result base.SliceResult[*entities.TrackingSession])
	Update(ctx context.Context, vc *entities.TrackingSession) error
}
type TrackingSessionRepository interface {
	base.Repository
	TrackingSessionCommon
}
