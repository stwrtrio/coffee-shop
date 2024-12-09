package repositories

import (
	"github.com/stwrtrio/coffee-shop/models"

	"gorm.io/gorm"
)

type EmailConfirmationRepository interface {
	CreateEmailConfirmation(user *models.Notification) error
}

type emailConfirmationRepository struct {
	db *gorm.DB
}

func NewEmailConfirmationRepository(db *gorm.DB) EmailConfirmationRepository {
	return &emailConfirmationRepository{db: db}
}

func (r *emailConfirmationRepository) CreateEmailConfirmation(user *models.Notification) error {
	return r.db.Create(user).Error
}
