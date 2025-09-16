package entities

import (
	"time"
)

type Customer struct {
	ID        int        `gorm:"primaryKey;autoIncrement"`
	Name      string     `gorm:"size:255;not null"`
	Email     string     `gorm:"size:255;not null"`
	Phone     string     `gorm:"size:14"`
	Password  string     `gorm:"column:password_hash"`
	Address   string     `gorm:"size:255"`
	IsActive  string     `gorm:"size:10"`
	CreatedAt time.Time  `gorm:"autoCreateTime"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime"`
	DeletedAt *time.Time `gorm:"index"`

	Vehicles []*Vehicle `gorm:"foreignKey:CreatedBy"`
}

func (Customer) TableName() string {
	return "customers"
}
