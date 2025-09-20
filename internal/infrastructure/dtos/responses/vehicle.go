package responses

type VehicleResponse struct {
	ID              string `json:"id"`
	PlateNo         string `json:"plateNo"`
	Model           string
	Brand           string
	Color           string
	ManufactureYear int
	MaxSpeed        int                `json:"max_speed"`
	FuelType        string             `json:"fuel_type"`
	IsActive        int                `json:"is_active"`
	Device          *DeviceGPSResponse `json:"device,omitempty"`
}

type VehicleShortInfo struct {
	ID      string `json:"id"`
	PlateNo string `json:"plate_no"`
	Model   string `json:"model"`
	Brand   string `json:"brand"`
}
