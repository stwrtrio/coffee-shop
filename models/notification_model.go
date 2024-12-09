package models

import "time"

type Notification struct {
	ID         string `gorm:"type:uuid;primaryKey"`
	CustomerID string `gorm:"type:uuid"`
	Email      string
	Code       string
	Type       string
	Status     string `gorm:"default:pending"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type EmailConfirmationMessage struct {
	CustomerID string `json:"customer_id"`
	Email      string `json:"email"`
	Code       string `json:"code"`
	Type       string `json:"type"`
}
