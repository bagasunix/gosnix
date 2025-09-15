package models

import (
	"time"

	"github.com/bagasunix/gosnix/internal/domain/customer"
)

type CustomerModel struct {
	ID        int
	Name      string
	Email     string
	Phone     string
	Password  string
	Address   string
	IsActive  string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	Vehicles []VehicleModel `gorm:"foreignKey:CreatedBy"`
}

func (c *CustomerModel) ToDomain() *customer.Customer {
	return &customer.Customer{
		ID:        c.ID,
		Name:      c.Name,
		Email:     c.Email,
		Phone:     c.Phone,
		Address:   c.Address,
		IsActive:  c.IsActive,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		DeletedAt: c.DeletedAt,
		Password:  c.Password,
	}
}
