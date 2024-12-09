package models

import "time"

type Customer struct {
	ID                      string    `gorm:"type:uuid;primaryKey"`
	FirstName               string    `gorm:"size:255;not null"`
	LastName                string    `gorm:"size:255;not null"`
	Email                   string    `gorm:"size:255;unique;not null"`
	PasswordHash            string    `gorm:"size:255;not null"`
	EmailConfirmationCode   string    `gorm:"size:255"`
	EmailConfirmationExpiry time.Time `gorm:"type:datetime;default:null"`
	IsEmailConfirmed        bool      `gorm:"default:false"`
	CreatedAt               time.Time
	UpdatedAt               time.Time
}

type RegisterRequest struct {
	FirstName string `json:"first_name" validate:"required,max=200"`
	LastName  string `json:"last_name" validate:"required,max=200"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type ConfirmCodeRequest struct {
	Email string `json:"email" validate:"required,email"`
	Code  string `json:"code" validate:"required"`
}
