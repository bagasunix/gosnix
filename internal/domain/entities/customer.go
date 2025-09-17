package entities

import (
	"time"
)

type Customer struct {
	ID             int        `gorm:"primaryKey;autoIncrement"`
	Name           string     `gorm:"size:255;not null"`
	Sex            int8       `json:"sex" gorm:"column:sex;not null"` // 1=male, 2=female
	DOB            *string    `json:"dob" gorm:"column:dob;size:10"`
	Email          string     `gorm:"size:255;not null"`
	Phone          string     `gorm:"size:14"`
	Password       string     `gorm:"column:password_hash"`
	Address        string     `gorm:"size:255"`
	Photo          string     `json:"photo" gorm:"column:photo;type:text"`
	CustomerStatus int8       `json:"customer_status" gorm:"column:customer_status;default:0"` // 1=active, 2=suspended, 0=inactive
	CreatedAt      time.Time  `gorm:"autoCreateTime"`
	UpdatedAt      time.Time  `gorm:"autoUpdateTime"`
	DeletedAt      *time.Time `gorm:"index"`

	Vehicles []*Vehicle `gorm:"foreignKey:CreatedBy"`
}

func (Customer) TableName() string {
	return "customers"
}
