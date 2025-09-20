package requests

import (
	"regexp"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type CreateCustomer struct {
	Name     string     `json:"name"`
	Sex      string     `json:"sex"`
	DOB      *time.Time `json:"dob"`
	Email    string     `json:"email"`
	Phone    string     `json:"phone"`
	Password string     `bson:"password"`
	Address  string     `json:"address"`
	Photo    string     `json:"photo"`

	DeviceGPS RegistrationGPS `json:"device_gps,omitempty"`

	Vehicle []CreateVehicle `json:"vehicle,omitempty"`
}

func (c CreateCustomer) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Email, validation.Required, is.Email),
		validation.Field(&c.Name, validation.Required),
		validation.Field(&c.Phone, validation.Length(0, 14), validation.Match(regexp.MustCompile(`^\d*$`)).Error("phone must be numeric")),
	)
}

type UpdateCustomer struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}

func (c UpdateCustomer) Validate() error {
	return validation.ValidateStruct(&c,
		// validation.Field(&c.Email, is.Email),
		validation.Field(&c.ID, validation.Match(regexp.MustCompile(`^\d*$`)).Error("id must be numeric")),
		validation.Field(&c.Phone, validation.Length(0, 14), validation.Match(regexp.MustCompile(`^\d*$`)).Error("phone must be numeric")),
	)
}
