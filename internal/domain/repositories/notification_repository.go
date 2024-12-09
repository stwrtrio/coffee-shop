package repositories

import (
	"github.com/stwrtrio/coffee-shop/models"

	"gorm.io/gorm"
)

type NotificationRepository interface {
	FetchPendingNotifications() ([]models.Notification, error)
	UpdateStatus(notificationID string, status string) error
}

type notificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) FetchPendingNotifications() ([]models.Notification, error) {
	var notifications []models.Notification
	err := r.db.Where("status = ?", "pending").Find(&notifications).Error
	return notifications, err
}

func (r *notificationRepository) UpdateStatus(notificationID string, status string) error {
	return r.db.Model(&models.Notification{}).
		Where("id = ?", notificationID).
		Update("status", status).Error
}
