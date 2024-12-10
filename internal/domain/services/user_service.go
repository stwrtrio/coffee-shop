package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/stwrtrio/coffee-shop/internal/domain/repositories"
	"github.com/stwrtrio/coffee-shop/models"
	"github.com/stwrtrio/coffee-shop/pkg/constants"
	"github.com/stwrtrio/coffee-shop/pkg/helpers"
	"github.com/stwrtrio/coffee-shop/pkg/kafka"
	"github.com/stwrtrio/coffee-shop/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	RegisterUser(ctx context.Context, input *models.RegisterRequest) error
	LoginUser(ctx context.Context, user *models.LoginRequest) (string, error)
	ConfirmCode(ctx context.Context, email, code string) error
}

type userService struct {
	config   *utils.Config
	userRepo repositories.UserRepository
	kafka    *kafka.KafkaClient
}

func NewUserService(config *utils.Config, userRepo repositories.UserRepository, kafka *kafka.KafkaClient) UserService {
	return &userService{config: config, userRepo: userRepo, kafka: kafka}
}

func (s *userService) RegisterUser(ctx context.Context, req *models.RegisterRequest) error {
	// Check if the email is already in use
	existingUser, err := s.userRepo.FindUserByEmail(ctx, req.Email)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return errors.New("email already in use")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}

	// Create user
	user := &models.User{
		ID:                      uuid.NewString(),
		Name:                    req.Name,
		Email:                   req.Email,
		Phone:                   req.Phone,
		Address:                 req.Address,
		PasswordHash:            string(hashedPassword),
		EmailConfirmationCode:   utils.GenerateConfirmationCode(),
		EmailConfirmationExpiry: time.Now().Add(15 * time.Minute),
		Role:                    "customer",
		IsEmailConfirmed:        false,
	}

	// Save the user to the database
	if err = s.userRepo.CreateUser(user); err != nil {
		return err
	}

	// Publish to Kafka
	message := map[string]string{
		"user_id": user.ID,
		"email":   user.Email,
		"code":    user.EmailConfirmationCode,
		"type":    string(constants.EmailTypeConfirmation),
	}
	messageBytes, _ := json.Marshal(message)

	err = s.kafka.Produce(s.config.Kafka.Topics.EmailConfirmation, messageBytes)
	if err != nil {
		log.Printf("Failed to send email notification: %v", err)
		return err
	}

	return nil
}

// Login checks if the user credentials are correct and generates a JWT token
func (s *userService) LoginUser(ctx context.Context, user *models.LoginRequest) (string, error) {
	// Fetch user from repository
	userResult, err := s.userRepo.FindUserByEmail(ctx, user.Email)
	if err != nil {
		return "", fmt.Errorf("user not found")
	}

	if !userResult.IsEmailConfirmed {
		return "", fmt.Errorf("user is not confirmed")
	}

	// Check if the password matches
	err = bcrypt.CompareHashAndPassword([]byte(userResult.PasswordHash), []byte(user.Password))
	if err != nil {
		return "", errors.New("invalid password")
	}

	// Generate JWT token
	token, err := helpers.GenerateJWTToken(&s.config.Jwt, userResult.ID, userResult.Email, userResult.Role)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %v", err)
	}

	return token, nil
}

func (s *userService) ConfirmCode(ctx context.Context, email, code string) error {
	// Fetch the user by email
	user, err := s.userRepo.FindUserByEmail(ctx, email)
	if err != nil {
		return err
	}

	// Check if the code matches and is not expired
	if user.EmailConfirmationCode != code {
		return errors.New("invalid confirmation code")
	}

	// Check if the code is expired
	if time.Now().After(user.EmailConfirmationExpiry) {
		return errors.New("confirmation code expired")
	}

	// Mark the email as confirmed
	user.IsEmailConfirmed = true
	user.EmailConfirmationCode = "" // Clear the confirmation code

	// Save the updated user data
	err = s.userRepo.Update(user)
	if err != nil {
		return err
	}

	return nil
}
