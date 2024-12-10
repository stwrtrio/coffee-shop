package repositories

import (
	"context"

	"github.com/stwrtrio/coffee-shop/models"
	"gorm.io/gorm"
)

type menuRepository struct {
	db *gorm.DB
}

type MenuRepository interface {
	CreateMenu(menu *models.Menu) error
	FindMenuByName(ctx context.Context, menuName string) (*models.Menu, error)
}

func NewMenuRepository(db *gorm.DB) MenuRepository {
	return &menuRepository{db: db}
}

func (r *menuRepository) CreateMenu(menu *models.Menu) error {
	return r.db.Create(menu).Error
}

func (r *menuRepository) FindMenuByName(ctx context.Context, menuName string) (*models.Menu, error) {
	var menu models.Menu
	err := r.db.WithContext(ctx).Where("name = ?", menuName).First(&menu).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &menu, err
}
