package services

import (
	"context"
	"errors"

	"github.com/stwrtrio/coffee-shop/internal/domain/repositories"
	"github.com/stwrtrio/coffee-shop/models"
)

type UserService interface {
	RegisterUser(ctx context.Context, input *models.User) error
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) RegisterUser(ctx context.Context, input *models.User) error {
	// Check if the email is already in use
	existingUser, err := s.userRepo.FindUserByEmail(ctx, input.Email)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return errors.New("email already in use")
	}

	// Save the user to the database
	return s.userRepo.CreateUser(input)
}
