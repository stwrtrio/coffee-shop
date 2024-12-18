package repositories

import (
	"context"

	"github.com/stwrtrio/coffee-shop/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	FindUserByID(ctx context.Context, userID string) (*models.User, error)
	FindUserByEmail(ctx context.Context, email string) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(ctx context.Context, user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindUserByID(ctx context.Context, userID string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("id = ?", userID).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &user, err
}

func (r *userRepository) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &user, err
}

func (r *userRepository) UpdateUser(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Updates(user).Where("id = ?", user.ID).Error
}
