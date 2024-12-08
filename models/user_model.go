package models

import "time"

type User struct {
	ID           string `gorm:"type:uuid;primaryKey"`
	FirstName    string `gorm:"size:255;not null"`
	LastName     string `gorm:"size:255;not null"`
	Email        string `gorm:"size:255;unique;not null"`
	PasswordHash string `gorm:"size:255;not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type RegisterRequest struct {
	FirstName string `json:"first_name" validate:"required,max=200"`
	LastName  string `json:"last_name" validate:"required,max=200"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=6"`
}
