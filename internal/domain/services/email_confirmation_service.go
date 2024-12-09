package services

import (
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/stwrtrio/coffee-shop/internal/domain/repositories"
	"github.com/stwrtrio/coffee-shop/models"
	"github.com/stwrtrio/coffee-shop/pkg/kafka"
	"github.com/stwrtrio/coffee-shop/pkg/utils"
)

type EmailConfirmationService interface {
	ConsumeEmailConfirmation() error
}

type emailConfirmationService struct {
	config                *utils.Config
	kafka                 *kafka.KafkaClient
	emailConfirmationRepo repositories.EmailConfirmationRepository
}

func NewEmailConfirmationRepository(config *utils.Config, kafka *kafka.KafkaClient, emailConfirmationRepo repositories.EmailConfirmationRepository) EmailConfirmationService {
	return &emailConfirmationService{
		config:                config,
		kafka:                 kafka,
		emailConfirmationRepo: emailConfirmationRepo,
	}
}

func (s *emailConfirmationService) ConsumeEmailConfirmation() error {
	// Subscribe to the Kafka topic for email confirmation
	err := s.kafka.Consumer.Subscribe(s.config.Kafka.Topics.EmailConfirmation, nil)
	if err != nil {
		log.Printf("Failed to subscribe to Kafka topic: %v", err)
		return err
	}

	for {
		// Read messages from Kafka
		msg, err := s.kafka.Consumer.ReadMessage(-1)
		if err != nil {
			log.Printf("Error reading message: %v", err)
			continue
		}

		// Unmarshal message into struct
		var emailMsg models.EmailConfirmationMessage
		err = json.Unmarshal(msg.Value, &emailMsg)
		if err != nil {
			log.Printf("Error unmarshalling message: %v", err)
			continue
		}

		// Insert into notification table
		notification := models.Notification{
			ID:        uuid.New().String(),
			UserID:    emailMsg.UserID,
			Email:     emailMsg.Email,
			Code:      emailMsg.Code,
			Type:      emailMsg.Type,
			Status:    "pending",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Insert notification into the database
		if err := s.emailConfirmationRepo.CreateEmailConfirmation(&notification); err != nil {
			log.Printf("Failed to insert notification: %v", err)
		} else {
			log.Println("Notification inserted successfully")
		}
	}
}
