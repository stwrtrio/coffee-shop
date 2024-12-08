package models

import "time"

type MenuItem struct {
	ID          string `gorm:"type:uuid;primaryKey"`
	Name        string `gorm:"size:255;not null"`
	Description string
	Price       float64 `gorm:"not null"`
	Category    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
