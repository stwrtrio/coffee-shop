package services

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/stwrtrio/coffee-shop/internal/domain/repositories"
	"github.com/stwrtrio/coffee-shop/models"
	"github.com/stwrtrio/coffee-shop/pkg/kafka"
	"github.com/stwrtrio/coffee-shop/pkg/utils"
)

type UserService interface {
	RegisterUser(ctx context.Context, input *models.User) error
}

type userService struct {
	userRepo repositories.UserRepository
	kafka    *kafka.KafkaClient
	config   *utils.Config
}

func NewUserService(userRepo repositories.UserRepository, config *utils.Config, kafka *kafka.KafkaClient) UserService {
	return &userService{userRepo: userRepo, config: config, kafka: kafka}
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
	if err = s.userRepo.CreateUser(input); err != nil {
		return err
	}

	// Publish to Kafka
	message := map[string]string{
		"user_id": input.ID,
		"email":   input.Email,
		"code":    input.EmailConfirmationCode,
		"type":    "email_confirmation",
	}
	messageBytes, _ := json.Marshal(message)

	err = s.kafka.Produce(s.config.Kafka.Topics.EmailConfirmation, messageBytes)
	if err != nil {
		log.Printf("Failed to send email notification: %v", err)
		return err
	}

	return nil
}
