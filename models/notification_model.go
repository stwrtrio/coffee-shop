package models

import "time"

type Notification struct {
	ID        string `gorm:"type:uuid;primaryKey"`
	UserID    string `gorm:"type:uuid"`
	Email     string
	Code      string
	Type      string
	Status    string `gorm:"default:pending"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type EmailConfirmationMessage struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Code   string `json:"code"`
	Type   string `json:"type"`
}
