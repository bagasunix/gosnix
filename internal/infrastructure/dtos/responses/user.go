package responses

import "time"

type UserResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Sex      string `json:"sex"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	IsActive string `json:"is_active"`

	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}
