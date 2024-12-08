package models

import "time"

type Inventory struct {
	ID        string `gorm:"type:uuid;primaryKey"`
	ItemName  string `gorm:"size:255;not null"`
	Quantity  int    `gorm:"not null"`
	Unit      string `gorm:"size:50;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
