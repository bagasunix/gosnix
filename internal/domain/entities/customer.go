package entities

import (
	"time"

	"gorm.io/gorm"
)

type Customer struct {
	ID           int            `gorm:"primaryKey;autoIncrement"`
	Name         string         `gorm:"size:255;not null"`
	Sex          int8           `gorm:"size:1;not null"` // Changed from Sex int to Gender string
	DOB          *string        `gorm:"type:date"`
	Email        string         `gorm:"size:255;unique;not null"`
	Phone        string         `gorm:"size:20;unique;not null"`
	PasswordHash string         `gorm:"size:255;not null"`
	Address      string         `gorm:"type:text"`
	Photo        string         `gorm:"type:text"`
	IsActive     int8           `gorm:"default:0;not null"`
	CreatedAt    time.Time      `gorm:"autoCreateTime"`
	UpdatedAt    *time.Time     `gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`

	Vehicles []Vehicle `gorm:"foreignKey:CustomerID;references:ID"`
}

func (Customer) TableName() string {
	return "customers"
}
