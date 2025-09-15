package location_update

import (
	"context"

	"github.com/gofrs/uuid"

	"github.com/bagasunix/gosnix/internal/domain/base"
)

type Common interface {
	Save(ctx context.Context, lu *LocationUpdate) error
	SaveBatch(ctx context.Context, updates []*LocationUpdate) error
	FindByVehicle(ctx context.Context, vehicleID uuid.UUID, limit, offset int, search string) (result base.SliceResult[*LocationUpdate])
	FindBySession(ctx context.Context, sessionID uuid.UUID, limit, offset int, search string) (result base.SliceResult[*LocationUpdate])
	// DeleteOlderThan(ctx context.Context, vehicleID uuid.UUID, before time.Time) error
}

type Repository interface {
	base.Repository
	Common
}
