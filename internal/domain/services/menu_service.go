package services

import (
	"github.com/stwrtrio/coffee-shop/internal/domain/repositories"
)

type MenuService struct {
	MenuRepo repositories.MenuRepository
}

func NewMenuService(repo repositories.MenuRepository) *MenuService {
	return &MenuService{MenuRepo: repo}
}
