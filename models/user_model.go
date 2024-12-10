package models

import "time"

type User struct {
	ID                      string    `gorm:"type:uuid;primaryKey" json:"id"`
	Name                    string    `gorm:"type:varchar(100)" json:"name"`
	Email                   string    `gorm:"type:varchar(255);unique;not null" json:"email"`
	Phone                   string    `gorm:"type:varchar(20)" json:"phone,omitempty"`
	Address                 string    `gorm:"type:text" json:"address,omitempty"`
	PasswordHash            string    `gorm:"type:text;not null" json:"-"`
	Role                    string    `gorm:"type:enum('customer', 'staff', 'admin');not null" json:"role"`
	EmailConfirmationCode   string    `gorm:"size:255"`
	EmailConfirmationExpiry time.Time `gorm:"type:datetime;default:null"`
	IsEmailConfirmed        bool      `gorm:"default:false"`
	CreatedAt               time.Time
	UpdatedAt               time.Time
}

type RegisterRequest struct {
	Name     string `json:"name" validate:"required,max=200"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type ConfirmCodeRequest struct {
	Email string `json:"email" validate:"required,email"`
	Code  string `json:"code" validate:"required"`
}
