package requests

import (
	"regexp"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gofrs/uuid"
)

type CreateVehicle struct {
	PlateNo         string           `json:"plate_no" example:"B 1234 ABC"`
	VIN             string           `json:"vin" example:"1HGCM82633A123456"`
	CategoryID      uuid.UUID        `json:"category_id" example:"3fa85f64-5717-4562-b3fc-2c963f66afa6"`
	Model           string           `json:"model" example:"Avanza"`
	Brand           string           `json:"brand" example:"Toyota"`
	Color           string           `json:"color" example:"Hitam"`
	ManufactureYear int              `json:"manufacture_year" example:"2020"`
	MaxSpeed        int              `json:"max_speed" example:"180"`
	FuelType        string           `json:"fuel_type" example:"Bensin"`
	DeviceGPS       *RegistrationGPS `json:"device_gps,omitempty"`
}

func (v CreateVehicle) Validate() error {
	currentYear := time.Now().Year()

	return validation.ValidateStruct(&v,
		validation.Field(&v.PlateNo,
			validation.Required,
			validation.Length(5, 12),
			validation.Match(regexp.MustCompile(`^[A-Z0-9\s]+$`)).Error("plate_no must be alphanumeric with spaces"),
		),
		validation.Field(&v.VIN,
			validation.Required,
			validation.Length(17, 17).Error("vin must be exactly 17 characters"),
			validation.Match(regexp.MustCompile(`^[A-HJ-NPR-Z0-9]+$`)).Error("vin must be alphanumeric (I,O,Q not allowed)"),
		),
		validation.Field(&v.Model, validation.Required),
		validation.Field(&v.Brand, validation.Required),
		validation.Field(&v.Color, validation.Required),
		validation.Field(&v.ManufactureYear,
			validation.Required,
			validation.Min(1900),
			validation.Max(currentYear),
		),
		validation.Field(&v.MaxSpeed,
			validation.Required,
			validation.Min(1).Error("max_speed must be greater than 0"),
		),
		validation.Field(&v.FuelType,
			validation.Required,
			validation.In("Bensin", "Solar", "Listrik", "Hybrid").Error("fuel_type must be Bensin, Solar, Listrik, or Hybrid"),
		),
	)
}

type BaseVehicle struct {
	CustomerID string `json:"customer_id"`
	BaseRequest
}

type UpdateVehicle struct {
	ID              string     `json:"id" example:"1"`
	PlateNo         *string    `json:"plate_no,omitempty" example:"B 1234 ABC"`
	VIN             *string    `json:"vin,omitempty" example:"1HGCM82633A123456"`
	CategoryID      *uuid.UUID `json:"category_id,omitempty" example:"3fa85f64-5717-4562-b3fc-2c963f66afa6"`
	Model           *string    `json:"model,omitempty" example:"Avanza"`
	Brand           *string    `json:"brand,omitempty" example:"Toyota"`
	Color           *string    `json:"color,omitempty" example:"Hitam"`
	ManufactureYear *int       `json:"manufacture_year,omitempty" example:"2020"`
	MaxSpeed        *int       `json:"max_speed,omitempty" example:"180"`
	FuelType        *string    `json:"fuel_type,omitempty" example:"Bensin"`
	Photo           *string    `json:"photo,omitempty" example:"https://example.com/photo.jpg"`
}

func (v UpdateVehicle) Validate() error {
	currentYear := time.Now().Year()

	return validation.ValidateStruct(&v,
		validation.Field(&v.ID,
			validation.Required,
			validation.Match(regexp.MustCompile(`^\d+$`)).Error("id must be numeric"),
		),
		validation.Field(&v.PlateNo,
			validation.NilOrNotEmpty,
			validation.Length(5, 12),
			validation.Match(regexp.MustCompile(`^[A-Z0-9\s]+$`)).Error("plate_no must be alphanumeric with spaces"),
		),
		validation.Field(&v.VIN,
			validation.NilOrNotEmpty,
			validation.Length(17, 17).Error("vin must be exactly 17 characters"),
			validation.Match(regexp.MustCompile(`^[A-HJ-NPR-Z0-9]+$`)).Error("vin must be alphanumeric (I,O,Q not allowed)"),
		),
		validation.Field(&v.ManufactureYear,
			validation.Min(1900),
			validation.Max(currentYear),
		),
		validation.Field(&v.MaxSpeed,
			validation.Min(1),
		),
		validation.Field(&v.FuelType,
			validation.In("Bensin", "Solar", "Listrik", "Hybrid").Error("fuel_type must be valid type"),
		),
	)
}
