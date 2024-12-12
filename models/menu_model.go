package models

import "time"

type Menu struct {
	ID              string    `gorm:"type:uuid;primaryKey" json:"id"`
	Name            string    `gorm:"type:varchar(255);not null" json:"name"`
	Description     string    `gorm:"type:text" json:"description"`
	Price           float64   `gorm:"type:decimal(10,2);not null" json:"price"`
	CategoryID      string    `gorm:"type:uuid;not null" json:"category_id"`
	Availability    bool      `gorm:"default:true" json:"availability"`
	ImageURL        string    `gorm:"type:varchar(255)" json:"image_url"`
	Ingredients     string    `gorm:"type:text" json:"ingredients"`
	PreparationTime int       `gorm:"type:int" json:"preparation_time"`
	Calories        int       `gorm:"type:int" json:"calories"`
	CreatedBy       string    `gorm:"type:uuid" json:"created_by"`
	UpdatedBy       string    `gorm:"type:uuid" json:"updated_by"`
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	IsDeleted       bool      `gorm:"default:false" json:"is_deleted"`
}

type MenuRequest struct {
	Name            string  `json:"name" validate:"required"`
	Description     string  `json:"description"`
	Price           float64 `json:"price"`
	CategoryID      string  `json:"category_id" validate:"required"`
	Availability    bool    `json:"availability"`
	ImageURL        string  `json:"image_url"`
	Ingredients     string  `json:"ingredients"`
	PreparationTime int     `json:"preparation_time"`
	Calories        int     `json:"calories"`

	MenuID    string
	CreatedBy string
	UpdatedBy string
}
