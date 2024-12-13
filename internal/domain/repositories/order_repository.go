package repositories

import (
	"context"

	"github.com/stwrtrio/coffee-shop/models"
	"gorm.io/gorm"
)

type orderRepository struct {
	db *gorm.DB
}

type OrderRepository interface {
	CreateMonthlyOrderTables(ctx context.Context, tableName string) error
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (s *orderRepository) CreateMonthlyOrderTables(ctx context.Context, tableName string) error {
	// Check if the table already exists
	if s.db.Migrator().HasTable(tableName) {
		return nil
	}
	return s.db.Table(tableName).Migrator().CreateTable(&models.Order{})
}
