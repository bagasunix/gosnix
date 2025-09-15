package requests

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"

	"github.com/bagasunix/gosnix/pkg/errors"
)

type BaseRequest struct {
	Search string `json:"search,omitempty"`
	Limit  string `json:"limit,omitempty"`
	Page   string `json:"page,omitempty"`
}

func (c BaseRequest) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Search, is.Alphanumeric),
		validation.Field(&c.Limit, is.UTFNumeric.Error("Limit harus berupa angka")),
		validation.Field(&c.Page, is.UTFNumeric.Error("Page harus berupa angka")),
	)
}

type EntityId struct {
	Id any `json:"id"`
}

func (v *EntityId) Validate() error {
	if validation.IsEmpty(v.Id) {
		return errors.ErrInvalidAttributes("id")
	}
	return nil
}
