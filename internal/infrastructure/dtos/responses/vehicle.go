package responses

type VehicleResponse struct {
	ID       string `json:"id"`
	PlateNo  string `json:"plateNo"`
	Model    string
	Brand    string
	Color    string
	Year     int
	MaxSpeed int    `json:"max_speed"`
	FuelType string `json:"fuel_type"`
	IsActive int    `json:"is_active"`
}
