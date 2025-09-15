package tracking_session

import (
	"context"

	"github.com/gofrs/uuid"

	"github.com/bagasunix/gosnix/internal/domain/base"
)

type Common interface {
	Save(ctx context.Context, vc *TrackingSession) error
	FindByParam(ctx context.Context, param map[string]any) (result base.SingleResult[*TrackingSession])
	FindByVehicle(ctx context.Context, vehicleID uuid.UUID) (result base.SliceResult[*TrackingSession])
	Update(ctx context.Context, vc *TrackingSession) error
}
type Repository interface {
	base.Repository
	Common
}
