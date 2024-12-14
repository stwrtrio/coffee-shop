package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/stwrtrio/coffee-shop/models"
	"gorm.io/gorm"
)

type orderRepository struct {
	db *gorm.DB
}

type OrderRepository interface {
	CreateMonthlyOrderTables(ctx context.Context, reqTime time.Time) error
	CreateOrder(tx *gorm.DB, order *models.Order, tableOrder string) error
	CreateOrderItem(tx *gorm.DB, orderItems *models.OrderItem, tableOrderItems string) error
	GetOrderWithItems(orderID string) (*models.Order, error)
	CreateOrderTransaction(order *models.Order) error
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) CreateMonthlyOrderTables(ctx context.Context, reqTime time.Time) error {
	var (
		order      *models.Order
		orderItems *models.OrderItem
	)

	// Check if the table already exists
	if r.db.WithContext(ctx).Migrator().HasTable(order.TableName(reqTime)) {
		return nil
	}
	if r.db.WithContext(ctx).Migrator().HasTable(orderItems.TableName(reqTime)) {
		return nil
	}

	// Dynamic Order Table
	if err := r.db.WithContext(ctx).Table(order.TableName(reqTime)).Migrator().CreateTable(&models.Order{}); err != nil {
		return err
	}

	// Dynamic Order Items Table
	if err := r.db.WithContext(ctx).Table(orderItems.TableName(reqTime)).Migrator().CreateTable(&models.OrderItem{}); err != nil {
		return err
	}

	return nil
}

func (r *orderRepository) GetOrderWithItems(orderID string) (*models.Order, error) {
	var order models.Order
	if err := r.db.Preload("OrderItems").First(&order, "id = ?", orderID).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) CreateOrder(tx *gorm.DB, order *models.Order, tableOrder string) error {
	return tx.Table(tableOrder).Omit("OrderItems").Create(order).Error
}

func (r *orderRepository) CreateOrderItem(tx *gorm.DB, orderItems *models.OrderItem, tableOrderItems string) error {
	return tx.Table(tableOrderItems).Create(orderItems).Error
}

func (r *orderRepository) CreateOrderTransaction(order *models.Order) error {
	var orderItems *models.OrderItem
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := r.CreateOrder(tx, order, order.TableName(order.CreatedAt)); err != nil {
			fmt.Println("err1:", err)
			return err
		}

		for _, orderItem := range order.OrderItems {
			if err := r.CreateOrderItem(tx, &orderItem, orderItems.TableName(order.CreatedAt)); err != nil {
				fmt.Println("err2:", err)
				return err
			}
		}

		return nil
	})
}
