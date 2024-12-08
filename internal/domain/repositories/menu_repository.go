package repositories

import (
	"gorm.io/gorm"
)

type menuRepository struct {
	db *gorm.DB
}

func NewMenuRepository(db *gorm.DB) MenuRepositoryImpl {
	return &menuRepository{db: db}
}
