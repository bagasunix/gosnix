package entities

import "time"

type User struct {
	ID       int    `gorm:"primaryKey;autoIncrement"`
	Name     string `gorm:"type:varchar(100);not null"`
	Sex      int    `gorm:"not null"`
	Email    string `gorm:"type:varchar(100);unique;not null"`
	Password string `gorm:"type:varchar(100);not null"`
	Role     string `gorm:"not null"`
	IsActive int

	CreatedAt *time.Time `gorm:"column:created_at;autoautoCreateTime"`
	UpdatedAt *time.Time `gorm:"column:created_at;autoUpdateTime"`
	DeletedAt *time.Time `gorm:"column:deleted_at"`
}

func (User) TableName() string {
	return "users"
}
