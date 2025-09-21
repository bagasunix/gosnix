package requests

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gofrs/uuid"
)

type CreateVehicle struct {
	PlateNo         string    `json:"plate_no" example:"B 1234 ABC"`
	VIN             string    `json:"vin" example:"1HGCM82633A123456"`
	CategoryID      uuid.UUID `json:"category_id" example:"3fa85f64-5717-4562-b3fc-2c963f66afa6"`
	Model           string    `json:"model" example:"Avanza"`
	Brand           string    `json:"brand" example:"Toyota"`
	Color           string    `json:"color" example:"Hitam"`
	ManufactureYear int       `json:"manufacture_year" example:"2020"`
	MaxSpeed        int       `json:"max_speed" example:"180"`
	FuelType        string    `json:"fuel_type" example:"Bensin"`

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
