package repositories

import (
	"gorm.io/gorm"
)

type menuRepository struct {
	db *gorm.DB
}

type MenuRepository interface {
}

func NewMenuRepository(db *gorm.DB) MenuRepository {
	return &menuRepository{db: db}
}
