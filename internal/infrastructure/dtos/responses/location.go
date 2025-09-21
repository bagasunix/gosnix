package responses

import "time"

type LocationResponse struct {
	Latitude  float64   `json:"latitude" example:"-6.200000"`
	Longitude float64   `json:"longitude" example:"106.816666"`
	Speed     float64   `json:"speed,omitempty" example:"60.5"`
	Heading   int       `json:"heading,omitempty" example:"180"`
	Timestamp time.Time `json:"timestamp" example:"2023-10-01T12:00:00Z"`
}
