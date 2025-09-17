package responses

import "time"

type CustomerResponse struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Sex            string `json:"sex"`
	Email          string `json:"email"`
	Phone          string `json:"phone"`
	Address        string
	Photo          string
	CustomerStatus string `json:"customer_status"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`

	Vehicle []VehicleResponse `json:"vehicle,omitempty"`
}
