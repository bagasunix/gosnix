package responses

import "time"

type CustomerResponse struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Sex      int8   `json:"sex"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Address  string
	Photo    string
	IsActive int8 `json:"is_active"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`

	Vehicle []VehicleResponse `json:"vehicle,omitempty"`
}
