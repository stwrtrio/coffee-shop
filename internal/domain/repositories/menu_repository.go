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
	CreateMenu(ctx context.Context, menu *models.Menu) error
	UpdateMenu(ctx context.Context, req *models.Menu) error
	GetAllMenus(ctx context.Context, offset, limit int) ([]models.Menu, error)
	GetMenuByID(ctx context.Context, menuID string) (*models.Menu, error)
	FindMenuByName(ctx context.Context, menuName string) (*models.Menu, error)
	DeleteMenu(ctx context.Context, menuID string) error
}

func NewMenuRepository(db *gorm.DB) MenuRepository {
	return &menuRepository{db: db}
}

func (r *menuRepository) CreateMenu(ctx context.Context, menu *models.Menu) error {
	return r.db.WithContext(ctx).Create(menu).Error
}

func (r *menuRepository) GetAllMenus(ctx context.Context, offset, limit int) ([]models.Menu, error) {
	var menus []models.Menu
	err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Where("is_deleted = ?", false).Find(&menus).Error
	return menus, err
}

func (r *menuRepository) FindMenuByName(ctx context.Context, menuName string) (*models.Menu, error) {
	var menu models.Menu
	err := r.db.WithContext(ctx).Where("name = ?", menuName).First(&menu).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &menu, err
}

func (r *menuRepository) GetMenuByID(ctx context.Context, menuID string) (*models.Menu, error) {
	var menu models.Menu
	err := r.db.WithContext(ctx).Where("id = ?", menuID).First(&menu).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return &menu, nil
}

func (r *menuRepository) UpdateMenu(ctx context.Context, req *models.Menu) error {
	return r.db.WithContext(ctx).Updates(req).Where("id = ?", req.ID).Error
}

func (r *menuRepository) DeleteMenu(ctx context.Context, menuID string) error {
	return r.db.WithContext(ctx).Model(&models.Menu{}).Where("id = ?", menuID).Update("is_deleted", true).Error
}
