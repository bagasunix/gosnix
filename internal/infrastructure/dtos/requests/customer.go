package requests

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type CreateCustomer struct {
	Name     string          `json:"name" example:"John Doe"`
	Sex      int8            `json:"sex" example:"1"`
	DOB      *string         `json:"dob" example:"1990-01-01"`
	Email    string          `json:"email" example:"demo1@gmail.com"`
	Phone    string          `json:"phone" example:"081234567890"`
	Password string          `json:"password" example:"password123"`
	Address  string          `json:"address" example:"Jl. Merdeka No. 123, Jakarta"`
	Photo    string          `json:"photo" example:"https://example.com/photo.jpg"`
	Vehicle  []CreateVehicle `json:"vehicle,omitempty"`
}

func (c CreateCustomer) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Name, validation.Required),
		validation.Field(&c.Email, validation.Required, is.Email),
		validation.Field(&c.Phone,
			validation.Required,
			validation.Length(10, 14),
			validation.Match(regexp.MustCompile(`^\d+$`)).Error("phone must be numeric"),
		),
		validation.Field(&c.Password, validation.Required, validation.Length(6, 50)),
		validation.Field(&c.DOB,
			validation.Required,
			validation.Date("2006-01-02").Error("dob must be in format YYYY-MM-DD"),
		),
	)
}

type UpdateCustomer struct {
	ID      string  `json:"-"`
	Name    *string `json:"name,omitempty" example:"John Doe"`
	Sex     *int8   `json:"sex,omitempty" example:"1"`
	DOB     *string `json:"dob,omitempty" example:"1990-01-01"`
	Email   *string `json:"email,omitempty" example:"demo1@gmail.com"`
	Phone   *string `json:"phone,omitempty" example:"081234567890"`
	Address *string `json:"address,omitempty" example:"Jl. Merdeka No. 123, Jakarta"`
	Photo   *string `json:"photo,omitempty" example:"https://example.com/photo.jpg"`
}

func (c UpdateCustomer) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.ID,
			validation.Required,
			validation.Match(regexp.MustCompile(`^\d+$`)).Error("id must be numeric"),
		),
		validation.Field(&c.DOB,
			validation.Required,
			validation.Date("2006-01-02").Error("dob must be in format YYYY-MM-DD"),
		),
	)
}
