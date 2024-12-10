package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/stwrtrio/coffee-shop/internal/domain/repositories"
	"github.com/stwrtrio/coffee-shop/models"
	"github.com/stwrtrio/coffee-shop/pkg/constants"
)

type MenuService interface {
	CreateMenu(ctx context.Context, req *models.MenuRequest) (*models.Menu, error)
}

type menuService struct {
	categoryRepo repositories.CategoryRepository
	menuRepo     repositories.MenuRepository
}

func NewMenuService(menuRepo repositories.MenuRepository, categoryRepo repositories.CategoryRepository) MenuService {
	return &menuService{menuRepo: menuRepo, categoryRepo: categoryRepo}
}

func (s *menuService) CreateMenu(ctx context.Context, req *models.MenuRequest) (*models.Menu, error) {
	var menu *models.Menu

	// Find Menu by name
	resultMenu, err := s.menuRepo.FindMenuByName(ctx, req.Name)
	if err != nil {
		return menu, err
	}

	// Return error if menu already exist
	if resultMenu != nil {
		return nil, errors.New(constants.ErrMenuAlreadyExists)
	}

	// Find category by category id
	category, err := s.categoryRepo.FindCategoryByID(ctx, req.CategoryID)
	if err != nil {
		return menu, err
	}

	// Return error if category not exist
	if category == nil {
		return nil, errors.New(constants.ErrCategoryNotExists)
	}

	menu = &models.Menu{
		ID:              uuid.NewString(),
		Name:            req.Name,
		Description:     req.Description,
		Price:           req.Price,
		CategoryID:      req.CategoryID,
		Availability:    req.Availability,
		ImageURL:        req.ImageURL,
		Ingredients:     req.Ingredients,
		PreparationTime: req.PreparationTime,
		Calories:        req.Calories,
		CreatedBy:       req.CreatedBy,
	}

	if err = s.menuRepo.CreateMenu(menu); err != nil {
		return nil, err
	}

	return menu, nil
}
