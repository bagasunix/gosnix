package vehicle_category

import (
	"context"

	"github.com/gofrs/uuid"

	"github.com/bagasunix/gosnix/internal/domain/base"
)

type Common interface {
	Save(ctx context.Context, vc *VehicleCategory) error
	FindByParam(ctx context.Context, param map[string]any) (result base.SingleResult[*VehicleCategory])
	FindAll(ctx context.Context, limit int, offset int, search string) (result base.SliceResult[*VehicleCategory])
	Update(ctx context.Context, vc *VehicleCategory) error
	Delete(ctx context.Context, id uuid.UUID) error
	CountVehicleCategory(ctx context.Context, search string) (int, error)
}
type Repository interface {
	base.Repository
	Common
}
