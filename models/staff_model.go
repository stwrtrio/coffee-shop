package models

import "time"

type Staff struct {
	ID        string `gorm:"type:uuid;primaryKey"`
	FirstName string `gorm:"size:255;not null"`
	LastName  string `gorm:"size:255;not null"`
	Role      string `gorm:"size:100;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type StaffShift struct {
	ID         string    `gorm:"type:uuid;primaryKey"`
	StaffID    string    `gorm:"type:uuid;not null"`
	ShiftStart time.Time `gorm:"not null"`
	ShiftEnd   time.Time `gorm:"not null"`
	CreatedAt  time.Time
	UpdatedAt  time.Time

	Staff Staff `gorm:"foreignKey:StaffID"`
}
