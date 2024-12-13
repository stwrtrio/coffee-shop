package models

import "time"

type Order struct {
	ID        string    `gorm:"type:varchar(36);primaryKey"`
	UserID    string    `gorm:"type:varchar(36);not null"`
	MenuID    string    `gorm:"type:varchar(36);not null"`
	Quantity  int       `gorm:"not null"`
	Total     float64   `gorm:"not null"`
	Status    string    `gorm:"type:enum('pending','completed','cancelled');default:'pending'"`
	CreatedBy string    `gorm:"type:varchar(36)" json:"created_by"`
	UpdatedBy string    `gorm:"type:varchar(36)" json:"updated_by"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	IsDeleted bool      `gorm:"default:false"`
}
