package requests

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
)

type GPSData struct {
	DeviceID  string
	Timestamp time.Time
	Lat       float64
	Lon       float64
	Speed     float64
	Direction float64
	Altitude  float64
	Battery   float64
	EventType string
}

type RegistrationGPS struct {
	Brand     string `json:"brand"`
	Model     string `json:"model"`
	Protocol  string `json:"protocol"`
	SecretKey string `bson:"secret_key"`
}

func (v RegistrationGPS) Validate() error {
	return validation.ValidateStruct(&v,
		validation.Field(&v.Brand, validation.Required),
		validation.Field(&v.Model, validation.Required),
		validation.Field(&v.Protocol, validation.Required),
		validation.Field(&v.SecretKey, validation.Required),
	)
}
