package responses

import "github.com/gofrs/uuid"

type VehicleResponse struct {
	ID              uuid.UUID          `json:"id" example:"3fa85f64-5717-4562-b3fc-2c963f66afa6"`
	PlateNo         string             `json:"plate_no" example:"B 1234 ABC"`
	Category        string             `json:"category" example:"SUV"`
	Model           string             `json:"model" example:"Avanza"`
	Brand           string             `json:"brand" example:"Toyota"`
	Color           string             `json:"color" example:"Hitam"`
	ManufactureYear int                `json:"manufacture_year" example:"2020"`
	MaxSpeed        int                `json:"max_speed" example:"180"`
	FuelType        string             `json:"fuel_type" example:"Bensin"`
	IsActive        int                `json:"is_active" example:"1"`
	Device          *DeviceGPSResponse `json:"device,omitempty"`
}

type VehicleShortInfo struct {
	ID      string `json:"id"`
	PlateNo string `json:"plate_no"`
	Model   string `json:"model"`
	Brand   string `json:"brand"`
}
