package services

import (
	"github.com/stwrtrio/coffee-shop/internal/domain/repositories"
	"github.com/stwrtrio/coffee-shop/models"
)

type MenuService interface {
	CreateMenu(menu *models.Menu) error
}

type menuService struct {
	menuRepo repositories.MenuRepository
}

func NewMenuService(repo repositories.MenuRepository) MenuService {
	return &menuService{menuRepo: repo}
}

func (s *menuService) CreateMenu(menu *models.Menu) error {
	return s.menuRepo.CreateMenu(menu)
}
