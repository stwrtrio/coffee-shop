package models

import "time"

type Notification struct {
	ID        string `gorm:"type:uuid;primaryKey"`
	UserID    string
	Email     string
	Code      string
	Type      string
	Status    string `gorm:"default:pending"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
