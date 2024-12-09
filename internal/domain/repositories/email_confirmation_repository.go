package repositories

import (
	"github.com/stwrtrio/coffee-shop/models"

	"gorm.io/gorm"
)

type EmailConfirmationRepository interface {
	CreateEmailConfirmation(customer *models.Notification) error
}

type emailConfirmationRepository struct {
	db *gorm.DB
}

func NewEmailConfirmationRepository(db *gorm.DB) EmailConfirmationRepository {
	return &emailConfirmationRepository{db: db}
}

func (r *emailConfirmationRepository) CreateEmailConfirmation(customer *models.Notification) error {
	return r.db.Create(customer).Error
}
