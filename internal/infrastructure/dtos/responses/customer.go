package responses

import "time"

type CustomerResponse struct {
	ID       int     `json:"id" example:"1"`
	Name     string  `json:"name" example:"John Doe"`
	DOB      *string `json:"dob,omitempty" example:"1990-01-01"`
	Sex      int8    `json:"sex" example:"1"`
	Email    string  `json:"email" example:"demo1@gmail.com"`
	Phone    string  `json:"phone" example:"+6281234567890"`
	Address  string  `json:"address" example:"Jl. Merdeka No. 123, Jakarta"`
	Photo    string  `json:"photo,omitempty" example:"https://example.com/photo.jpg"`
	IsActive int8    `json:"is_active" example:"1"`

	CreatedAt time.Time  `json:"created_at" example:"2023-10-01T12:00:00Z"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" example:"2023-10-02T12:00:00Z"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" example:"2023-10-03T12:00:00Z,omitempty"`

	Vehicle []VehicleResponse `json:"vehicle,omitempty"`
}
