package scheduler

import (
	"log"
	"time"

	"github.com/stwrtrio/coffee-shop/internal/domain/repositories"
	"github.com/stwrtrio/coffee-shop/models"
	"github.com/stwrtrio/coffee-shop/pkg/email"
)

type Scheduler interface {
	StartScheduler()
}

type scheduler struct {
	NotificationRepo repositories.NotificationRepository
	OrderRepo        repositories.OrderRepository
	EmailService     email.EmailService
}

func NewScheduler(
	notificationRepo repositories.NotificationRepository,
	OrderRepo repositories.OrderRepository,
	emailService email.EmailService) Scheduler {
	return &scheduler{
		NotificationRepo: notificationRepo,
		OrderRepo:        OrderRepo,
		EmailService:     emailService,
	}
}

func (s *scheduler) StartScheduler() {
	notificationChan := make(chan models.Notification, 10) // Buffered channel for notifications
	doneChan := make(chan bool)
	shutdownChan := make(chan struct{}) // For graceful shutdown

	// Start workers to process notifications
	for i := 0; i < 5; i++ {
		go s.processNotifications(notificationChan, doneChan)
	}

	// Start the ticker to fetch notifications periodically
	ticker := time.NewTicker(1 * time.Minute) // Run every 1 minute
	defer ticker.Stop()

	orderTableTicker := time.NewTicker(1 * time.Minute) // Run daily
	defer orderTableTicker.Stop()

	go func() {
		// Gracefully shutdown the scheduler when done
		<-doneChan
		close(shutdownChan) // Notify that the scheduler is shutting down
		log.Println("Scheduler shutting down gracefully.")
	}()

	for {
		select {
		case <-ticker.C:
			// Fetch pending notifications
			notifications, err := s.NotificationRepo.FetchPendingNotifications()
			if err != nil {
				log.Printf("Error fetching pending notifications: %v", err)
				continue
			}

			// Send notifications to the channel
			for _, notification := range notifications {
				notificationChan <- notification
			}

		case <-orderTableTicker.C:
			// Create monthly order tables
			err := s.createMonthlyTables()
			if err != nil {
				log.Printf("Error creating monthly order tables: %v", err)
			}

		case <-shutdownChan:
			// Shut down scheduler gracefully
			log.Println("Scheduler is stopping.")
			return
		}
	}
}
