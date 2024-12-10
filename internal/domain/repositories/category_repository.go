package repositories

import (
	"context"

	"github.com/stwrtrio/coffee-shop/models"
	"gorm.io/gorm"
)

type categoryRepository struct {
	db *gorm.DB
}

type CategoryRepository interface {
	CreateCategory(category *models.Categories) error
	GetAllCategories(ctx context.Context) ([]models.Categories, error)
	FindCategoryByID(ctx context.Context, categoryID string) (*models.Categories, error)
	FindCategoryByName(ctx context.Context, categoryName string) (*models.Categories, error)
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) CreateCategory(category *models.Categories) error {
	return r.db.Create(category).Error
}

func (r *categoryRepository) GetAllCategories(ctx context.Context) ([]models.Categories, error) {
	var categories []models.Categories
	err := r.db.WithContext(ctx).Where("is_deleted = ?", false).Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *categoryRepository) FindCategoryByID(ctx context.Context, categoryID string) (*models.Categories, error) {
	var category models.Categories
	err := r.db.WithContext(ctx).Where("id = ?", categoryID).First(&category).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &category, err
}

func (r *categoryRepository) FindCategoryByName(ctx context.Context, categoryName string) (*models.Categories, error) {
	var category models.Categories
	err := r.db.WithContext(ctx).Where("name = ?", categoryName).First(&category).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &category, err
}
