package responses

import "time"

type LocationResponse struct {
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Speed     float64   `json:"speed,omitempty"`
	Heading   int       `json:"heading,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}
