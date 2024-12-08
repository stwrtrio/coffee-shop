package models

import "time"

type Payment struct {
	ID            string  `gorm:"type:uuid;primaryKey"`
	OrderID       string  `gorm:"type:uuid;not null"`
	PaymentMethod string  `gorm:"size:50;not null"`
	Amount        float64 `gorm:"not null"`
	Status        string  `gorm:"default:'Completed'"`
	TransactionID string  `gorm:"size:255"`
	CreatedAt     time.Time
	UpdatedAt     time.Time

	Order Order `gorm:"foreignKey:OrderID"`
}
