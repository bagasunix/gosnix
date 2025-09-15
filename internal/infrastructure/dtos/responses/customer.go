package responses

import "time"

type CustomerResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Address  string
	IsActive string `json:"is_active"`

	CreatedBy int        `json:"created_by"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`

	Vehicle []VehicleResponse `json:"vehicle,omitempty"`
}
