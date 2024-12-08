package models

import "time"

type Order struct {
	ID         string  `gorm:"type:uuid;primaryKey"`
	CustomerID string  `gorm:"type:uuid;not null"`
	TotalPrice float64 `gorm:"not null"`
	Status     string  `gorm:"default:'Pending'"`
	CreatedAt  time.Time
	UpdatedAt  time.Time

	Customer Customer `gorm:"foreignKey:CustomerID"`
}

type OrderItem struct {
	ID         string  `gorm:"type:uuid;primaryKey"`
	OrderID    string  `gorm:"type:uuid;not null"`
	MenuItemID string  `gorm:"type:uuid;not null"`
	Quantity   int     `gorm:"not null"`
	UnitPrice  float64 `gorm:"not null"`
	CreatedAt  time.Time
	UpdatedAt  time.Time

	Order    Order    `gorm:"foreignKey:OrderID"`
	MenuItem MenuItem `gorm:"foreignKey:MenuItemID"`
}
