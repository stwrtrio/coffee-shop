package services

import (
	"github.com/stwrtrio/coffee-shop/internal/domain/repositories"
)

type MenuService struct {
	MenuRepo repositories.MenuRepositoryImpl
}

func NewMenuService(repo repositories.MenuRepositoryImpl) *MenuService {
	return &MenuService{MenuRepo: repo}
}
