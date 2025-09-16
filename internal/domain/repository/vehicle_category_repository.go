package repository

import (
	"context"

	"github.com/gofrs/uuid"

	"github.com/bagasunix/gosnix/internal/domain/base"
	"github.com/bagasunix/gosnix/internal/domain/entities"
)

type VehicleCategoryCommon interface {
	Save(ctx context.Context, vc *entities.VehicleCategory) error
	FindByParam(ctx context.Context, param map[string]any) (result base.SingleResult[*entities.VehicleCategory])
	FindAll(ctx context.Context, limit int, offset int, search string) (result base.SliceResult[*entities.VehicleCategory])
	Update(ctx context.Context, vc *entities.VehicleCategory) error
	Delete(ctx context.Context, id uuid.UUID) error
	CountVehicleCategory(ctx context.Context, search string) (int, error)
}
type VehicleCategoryRepository interface {
	base.Repository
	VehicleCategoryCommon
}
