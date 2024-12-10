package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/stwrtrio/coffee-shop/internal/domain/repositories"
	"github.com/stwrtrio/coffee-shop/models"
	"gorm.io/gorm"
)

type CategoryService interface {
	CreateCategory(ctx context.Context, req *models.CategoryRequest) (*models.Categories, error)
}

type categoryService struct {
	categoryRepo repositories.CategoryRepository
}

func NewCategoryService(repo repositories.CategoryRepository) CategoryService {
	return &categoryService{categoryRepo: repo}
}

func (s *categoryService) CreateCategory(ctx context.Context, req *models.CategoryRequest) (*models.Categories, error) {
	var categoryResult *models.Categories
	_, err := s.categoryRepo.FindCategoryByName(ctx, req.Name)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return categoryResult, err
		}
	}

	category := &models.Categories{
		ID:          uuid.NewString(),
		Name:        req.Name,
		Description: req.Description,
	}

	if err = s.categoryRepo.CreateCategory(category); err != nil {
		return categoryResult, err
	}

	categoryResult, err = s.categoryRepo.FindCategoryByID(ctx, category.ID)
	if err != nil {
		return categoryResult, err
	}

	return categoryResult, nil
}
