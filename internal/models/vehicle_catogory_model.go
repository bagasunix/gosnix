package models

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/bagasunix/gosnix/internal/domain/vehicle_category"
)

type VehicleCategoryModel struct {
	ID          uuid.UUID
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time

	Vehicles []VehicleModel `gorm:"foreignKey:CategoryID"`
}

func (vc *VehicleCategoryModel) ToDomain() *vehicle_category.VehicleCategory {
	return &vehicle_category.VehicleCategory{
		ID:          vc.ID,
		Name:        vc.Name,
		Description: vc.Description,
		CreatedAt:   vc.CreatedAt,
		UpdatedAt:   vc.UpdatedAt,
		DeletedAt:   vc.DeletedAt,
	}
}
