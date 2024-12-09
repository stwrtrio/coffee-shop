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
	"github.com/stwrtrio/coffee-shop/pkg/jwt"
	"github.com/stwrtrio/coffee-shop/pkg/kafka"
	"github.com/stwrtrio/coffee-shop/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

type CustomerService interface {
	RegisterCustomer(ctx context.Context, input *models.RegisterRequest) error
	LoginCustomer(ctx context.Context, customer *models.LoginRequest) (string, error)
	ConfirmCode(ctx context.Context, email, code string) error
}

type customerService struct {
	config       *utils.Config
	customerRepo repositories.CustomerRepository
	kafka        *kafka.KafkaClient
}

func NewCustomerService(config *utils.Config, customerRepo repositories.CustomerRepository, kafka *kafka.KafkaClient) CustomerService {
	return &customerService{config: config, customerRepo: customerRepo, kafka: kafka}
}

func (s *customerService) RegisterCustomer(ctx context.Context, req *models.RegisterRequest) error {
	// Check if the email is already in use
	existingCustomer, err := s.customerRepo.FindCustomerByEmail(ctx, req.Email)
	if err != nil {
		return err
	}
	if existingCustomer != nil {
		return errors.New("email already in use")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}

	// Create customer
	customer := &models.Customer{
		ID:                      uuid.NewString(),
		FirstName:               req.FirstName,
		LastName:                req.LastName,
		Email:                   req.Email,
		PasswordHash:            string(hashedPassword),
		EmailConfirmationCode:   utils.GenerateConfirmationCode(),
		EmailConfirmationExpiry: time.Now().Add(15 * time.Minute),
		IsEmailConfirmed:        false,
	}

	// Save the customer to the database
	if err = s.customerRepo.CreateCustomer(customer); err != nil {
		return err
	}

	// Publish to Kafka
	message := map[string]string{
		"customer_id": customer.ID,
		"email":       customer.Email,
		"code":        customer.EmailConfirmationCode,
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

// Login checks if the customer credentials are correct and generates a JWT token
func (s *customerService) LoginCustomer(ctx context.Context, customer *models.LoginRequest) (string, error) {
	// Fetch customer from repository
	customerResult, err := s.customerRepo.FindCustomerByEmail(ctx, customer.Email)
	if err != nil {
		return "", fmt.Errorf("customer not found")
	}

	if !customerResult.IsEmailConfirmed {
		return "", fmt.Errorf("customer is not confirmed")
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

func (s *customerService) ConfirmCode(ctx context.Context, email, code string) error {
	// Fetch the customer by email
	customer, err := s.customerRepo.FindCustomerByEmail(ctx, email)
	if err != nil {
		return err
	}

	// Check if the code matches and is not expired
	if customer.EmailConfirmationCode != code {
		return errors.New("invalid confirmation code")
	}

	// Check if the code is expired
	if time.Now().After(customer.EmailConfirmationExpiry) {
		return errors.New("confirmation code expired")
	}

	// Mark the email as confirmed
	customer.IsEmailConfirmed = true
	customer.EmailConfirmationCode = "" // Clear the confirmation code

	// Save the updated customer data
	err = s.customerRepo.Update(customer)
	if err != nil {
		return err
	}

	return nil
}
