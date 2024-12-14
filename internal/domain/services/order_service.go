package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/stwrtrio/coffee-shop/internal/domain/repositories"
	"github.com/stwrtrio/coffee-shop/models"
	"github.com/stwrtrio/coffee-shop/pkg/constants"
)

type OrderService interface {
	CreateOrder(ctx context.Context, req *models.RequestOrder, userID string) error
}

type orderService struct {
	orderRepo repositories.OrderRepository
	menuRepo  repositories.MenuRepository
}

func NewOrderService(orderRepo repositories.OrderRepository, menuRepo repositories.MenuRepository) OrderService {
	return &orderService{orderRepo: orderRepo, menuRepo: menuRepo}
}

func (s *orderService) CreateOrder(ctx context.Context, req *models.RequestOrder, staffID string) error {
	// Calculate total
	var total float64

	// Create Order
	order := &models.Order{
		ID:        uuid.New().String(),
		UserID:    req.UserID,
		Status:    constants.OrderStatusPending,
		CreatedBy: staffID, // Staff ID from context
		CreatedAt: req.CreatedAt,
	}

	orderItems := []models.OrderItem{}
	for _, item := range req.Items {
		menu, err := s.menuRepo.GetMenuByID(ctx, item.MenuID)
		if err != nil {
			return err
		}

		if menu == nil {
			return errors.New("menu id not found")
		}

		total += float64(item.Quantity) * menu.Price

		orderItems = append(orderItems, models.OrderItem{
			ID:       uuid.New().String(),
			OrderID:  order.ID,
			MenuID:   item.MenuID,
			Quantity: item.Quantity,
			Price:    menu.Price,
		})
	}

	order.OrderItems = orderItems
	order.Total = total

	if err := s.orderRepo.CreateOrderTransaction(order); err != nil {
		return err
	}

	return nil
}
