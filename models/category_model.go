package models

import "time"

type Categories struct {
	ID          string    `gorm:"type:uuid;primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(100);unique;not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type CategoryRequest struct {
	Name        string `json:"name" validate:"required,max=200"`
	Description string `json:"description"`
}
