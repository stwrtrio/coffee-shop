package repositories

import (
	"github.com/stwrtrio/coffee-shop/models"
	"gorm.io/gorm"
)

type menuRepository struct {
	db *gorm.DB
}

type MenuRepository interface {
	CreateMenu(menu *models.Menu) error
}

func NewMenuRepository(db *gorm.DB) MenuRepository {
	return &menuRepository{db: db}
}

func (r *menuRepository) CreateMenu(menu *models.Menu) error {
	return r.db.Create(menu).Error
}
