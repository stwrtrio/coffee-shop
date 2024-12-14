package models

import (
	"fmt"
	"time"
)

type RequestOrder struct {
	UserID    string              `json:"user_id" validate:"required"`
	Items     []RequestOrderItems `json:"items" validate:"required,dive"`
	CreatedAt time.Time           `json:"created_at" validate:"required"`
}

type RequestOrderItems struct {
	MenuID   string `json:"menu_id" validate:"required"`
	Quantity int    `json:"quantity" validate:"required,min=1"`
}

type Order struct {
	ID         string      `gorm:"type:varchar(36);primaryKey"`
	UserID     string      `gorm:"type:varchar(36);not null"`
	Total      float64     `gorm:"not null"`
	Status     string      `gorm:"type:enum('pending','completed','cancelled');default:'pending'"`
	CreatedBy  string      `gorm:"type:varchar(36)" json:"created_by"`
	UpdatedBy  string      `gorm:"type:varchar(36)" json:"updated_by"`
	CreatedAt  time.Time   `gorm:"autoCreateTime"`
	UpdatedAt  time.Time   `gorm:"autoUpdateTime"`
	IsDeleted  bool        `gorm:"default:false"`
	OrderItems []OrderItem `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE"`
}

type OrderItem struct {
	ID        string    `gorm:"type:varchar(36);primaryKey"`
	OrderID   string    `gorm:"type:varchar(36);not null;index"`
	MenuID    string    `gorm:"type:varchar(36);not null"`
	Quantity  int       `gorm:"not null"`
	Price     float64   `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (o *Order) TableName(reqTime time.Time) string {
	yearMonth := reqTime.Format("200601") // Format as "YYYYMM"
	return fmt.Sprintf("orders_%s", yearMonth)
}

func (oi *OrderItem) TableName(reqTime time.Time) string {
	yearMonth := reqTime.Format("200601") // Format as "YYYYMM"
	return fmt.Sprintf("order_items_%s", yearMonth)
}
