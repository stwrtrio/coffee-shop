package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/stwrtrio/coffee-shop/internal/domain/repositories"
	"github.com/stwrtrio/coffee-shop/models"
	"github.com/stwrtrio/coffee-shop/pkg/constants"
	"github.com/stwrtrio/coffee-shop/pkg/utils"

	"github.com/go-redis/redis/v8"
)

type MenuService interface {
	CreateMenu(ctx context.Context, req *models.MenuRequest) (*models.Menu, error)
	GetAllMenus(ctx context.Context, page, limit int, useCache string) ([]models.Menu, error)
	GetMenusFromCache(ctx context.Context, page, limit int) ([]models.Menu, error)
}

type menuService struct {
	config       *utils.Config
	redis        *redis.Client
	categoryRepo repositories.CategoryRepository
	menuRepo     repositories.MenuRepository
}

func NewMenuService(config *utils.Config, redis *redis.Client, menuRepo repositories.MenuRepository, categoryRepo repositories.CategoryRepository) MenuService {
	return &menuService{config: config, redis: redis, menuRepo: menuRepo, categoryRepo: categoryRepo}
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

func (s *menuService) GetAllMenus(ctx context.Context, page, limit int, useCache string) ([]models.Menu, error) {
	var menus []models.Menu

	if useCache == "true" || useCache == "1" {
		menus, err := s.GetMenusFromCache(ctx, page, limit)
		if len(menus) > 0 {
			return menus, nil
		} else if err != nil {
			return menus, err
		}
	}

	// Cache miss - fetch from database
	offset := (page - 1) * limit
	menus, err := s.menuRepo.GetAllMenus(ctx, offset, limit)
	if err != nil {
		return nil, err
	}

	// Store result in Redis cache
	cacheKey := fmt.Sprintf(constants.MenusCacheKey, page, limit)
	cacheData, _ := json.Marshal(menus)
	s.redis.Set(ctx, cacheKey, cacheData, 5*time.Minute)

	return menus, nil
}

// GetMenusFromCache fetches menus from Redis cache.
func (s *menuService) GetMenusFromCache(ctx context.Context, page, limit int) ([]models.Menu, error) {
	var cacheMenus []models.Menu

	// Define cache key
	cacheKey := fmt.Sprintf(constants.MenusCacheKey, page, limit)

	// Get cached data
	cachedData, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		// Cache hit - Unmarshal data
		if err := json.Unmarshal([]byte(cachedData), &cacheMenus); err != nil {
			return nil, fmt.Errorf("failed to unmarshal cached data: %v", err)
		}
		return cacheMenus, nil
	}

	return cacheMenus, nil
}
