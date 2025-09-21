package requests

import (
	"regexp"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type CreateCustomer struct {
	Name     string     `json:"name" example:"John Doe"`
	Sex      string     `json:"sex" example:"1"`
	DOB      *time.Time `json:"dob" example:"1990-01-01"`
	Email    string     `json:"email" example:"demo1@gmail.com"`
	Phone    string     `json:"phone" example:"081234567890"`
	Password string     `json:"password" example:"password123"`
	Address  string     `json:"address" example:"Jl. Merdeka No. 123, Jakarta"`
	Photo    string     `json:"photo" example:"https://example.com/photo.jpg"`

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
