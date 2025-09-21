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
	IMEI      string `json:"imei" example:"123456789012345"`
	Brand     string `json:"brand" example:"Xiaomi"`
	Model     string `json:"model" example:"Mi Band 5"`
	Protocol  string `json:"protocol" example:"TCP"`
	SecretKey string `bson:"secret_key" example:"s3cr3tK3y"`
}

func (v RegistrationGPS) Validate() error {
	return validation.ValidateStruct(&v,
		validation.Field(&v.Brand, validation.Required),
		validation.Field(&v.Model, validation.Required),
		validation.Field(&v.Protocol, validation.Required),
		validation.Field(&v.SecretKey, validation.Required),
	)
}
