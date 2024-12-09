package services

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/stwrtrio/coffee-shop/internal/domain/repositories"
	"github.com/stwrtrio/coffee-shop/models"
	"github.com/stwrtrio/coffee-shop/pkg/constants"
	"github.com/stwrtrio/coffee-shop/pkg/kafka"
	"github.com/stwrtrio/coffee-shop/pkg/utils"
)

type CustomerService interface {
	RegisterCustomer(ctx context.Context, input *models.Customer) error
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
