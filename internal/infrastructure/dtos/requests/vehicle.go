package requests

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gofrs/uuid"
)

type CreateVehicle struct {
	PlateNo         string    `json:"plate_no"`
	VIN             string    `json:"vin"`
	CategoryID      uuid.UUID `json:"category_id"`
	Model           string    `json:"model"`
	Brand           string    `json:"brand"`
	Color           string    `json:"color"`
	ManufactureYear int       `json:"manufacture_year"`
	MaxSpeed        int       `json:"max_speed"`
	FuelType        string    `json:"fuel_type"`

	DeviceGPS RegistrationGPS `json:"device_gps,omitempty"`
}

func (v CreateVehicle) Validate() error {
	return validation.ValidateStruct(&v,
		validation.Field(&v.PlateNo, validation.Required),
		validation.Field(&v.Model, validation.Required),
		validation.Field(&v.Brand, validation.Required),
		validation.Field(&v.Color, validation.Required),
		validation.Field(&v.ManufactureYear, validation.Required),
		validation.Field(&v.FuelType, validation.Required),
	)
}

type BaseVehicle struct {
	CustomerID string `json:"customer_id"`
	BaseRequest
}
