package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/stwrtrio/coffee-shop/internal/domain/repositories"
	"github.com/stwrtrio/coffee-shop/models"
	"github.com/stwrtrio/coffee-shop/pkg/constants"
	"github.com/stwrtrio/coffee-shop/pkg/jwt"
	"github.com/stwrtrio/coffee-shop/pkg/kafka"
	"github.com/stwrtrio/coffee-shop/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

type CustomerService interface {
	RegisterCustomer(ctx context.Context, input *models.Customer) error
	LoginCustomer(ctx context.Context, customer *models.LoginRequest) (string, error)
}

type customerService struct {
	config       *utils.Config
	customerRepo repositories.CustomerRepository
	kafka        *kafka.KafkaClient
}

func NewCustomerService(config *utils.Config, customerRepo repositories.CustomerRepository, kafka *kafka.KafkaClient) CustomerService {
	return &customerService{config: config, customerRepo: customerRepo, kafka: kafka}
}

func (s *customerService) RegisterCustomer(ctx context.Context, input *models.Customer) error {
	// Check if the email is already in use
	existingCustomer, err := s.customerRepo.FindCustomerByEmail(ctx, input.Email)
	if err != nil {
		return err
	}
	if existingCustomer != nil {
		return errors.New("email already in use")
	}

	// Save the customer to the database
	if err = s.customerRepo.CreateCustomer(input); err != nil {
		return err
	}

	// Publish to Kafka
	message := map[string]string{
		"customer_id": input.ID,
		"email":       input.Email,
		"code":        input.EmailConfirmationCode,
		"type":        string(constants.EmailTypeConfirmation),
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
func (s *customerService) LoginCustomer(ctx context.Context, customer *models.LoginRequest) (string, error) {
	// Fetch customer from repository
	customerResult, err := s.customerRepo.FindCustomerByEmail(ctx, customer.Email)
	if err != nil {
		return "", errors.New("user not found")
	}

	// Check if the password matches
	err = bcrypt.CompareHashAndPassword([]byte(customerResult.PasswordHash), []byte(customer.Password))
	if err != nil {
		return "", errors.New("invalid password")
	}

	// Generate JWT token
	token, err := jwt.GenerateJWTToken(s.config, customerResult.ID, customerResult.Email)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %v", err)
	}

	return token, nil
}
