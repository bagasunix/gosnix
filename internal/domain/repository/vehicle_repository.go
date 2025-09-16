package repository

import (
	"context"

	"github.com/gofrs/uuid"

	"github.com/bagasunix/gosnix/internal/domain/base"
	"github.com/bagasunix/gosnix/internal/domain/entities"
)

type VehicleCommon interface {
	Save(ctx context.Context, v *entities.Vehicle) error
	FindByCustomer(ctx context.Context, customerID uuid.UUID) (result base.SliceResult[[]*entities.Vehicle])
	FindByParam(ctx context.Context, params map[string]interface{}) (result base.SingleResult[*entities.Vehicle])
	Update(ctx context.Context, v *entities.Vehicle) error
	Delete(ctx context.Context, id uuid.UUID) error
}
type VehicleRepository interface {
	base.Repository
	VehicleCommon
}
