package entities

import "time"

type AuditLog struct {
	ID         int       `gorm:"primaryKey;autoIncrement"`
	UserID     string    `gorm:"column:user_id"`
	Method     string    `gorm:"type:varchar(50);not null"`
	Endpoint   string    `gorm:"column:url;not null"`
	StatusCode string    `gorm:"column:status_code"`
	Request    string    `gorm:"column:request"`
	Response   string    `gorm:"column:response"`
	IPAddress  string    `gorm:"column:ip_address;type:varchar(45);not null"`
	UserAgent  string    `gorm:"column:user_agent;type:text;not null"`
	CreatedAt  time.Time `gorm:"column:created_at;type:timestamp;default:current_timestamp"`
}

// TableName sets the insert table name for this struct type
func (AuditLog) TableName() string {
	return "audit_logs"
}
