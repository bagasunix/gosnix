package repository

import (
	"context"

	"github.com/gofrs/uuid"

	"github.com/bagasunix/gosnix/internal/domain/base"
	"github.com/bagasunix/gosnix/internal/domain/entities"
)

type LocationUpdateCommon interface {
	Save(ctx context.Context, lu *entities.LocationUpdate) error
	SaveBatch(ctx context.Context, updates []*entities.LocationUpdate) error
	FindByVehicle(ctx context.Context, vehicleID uuid.UUID, limit, offset int, search string) (result base.SliceResult[*entities.LocationUpdate])
	FindBySession(ctx context.Context, sessionID uuid.UUID, limit, offset int, search string) (result base.SliceResult[*entities.LocationUpdate])
	// DeleteOlderThan(ctx context.Context, vehicleID uuid.UUID, before time.Time) error
}

type LocationUpdateRepository interface {
	base.Repository
	LocationUpdateCommon
}
