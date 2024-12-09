package scheduler

import (
	"fmt"
	"log"
	"time"

	"github.com/stwrtrio/coffee-shop/internal/domain/repositories"
	"github.com/stwrtrio/coffee-shop/models"
	"github.com/stwrtrio/coffee-shop/pkg/constants"
	"github.com/stwrtrio/coffee-shop/pkg/email"
)

type NotificationScheduler struct {
	NotificationRepo repositories.NotificationRepository
	EmailService     email.EmailService
}

func NewNotificationScheduler(notificationRepo repositories.NotificationRepository, emailService email.EmailService) *NotificationScheduler {
	return &NotificationScheduler{
		NotificationRepo: notificationRepo,
		EmailService:     emailService,
	}
}

func (s *NotificationScheduler) StartScheduler() {
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

		case <-shutdownChan:
			// Shut down scheduler gracefully
			log.Println("Scheduler is stopping.")
			return
		}
	}
}

func (s *NotificationScheduler) processNotifications(notificationChan <-chan models.Notification, doneChan chan<- bool) {
	for notification := range notificationChan {
		subject := constants.NotificationConstants.EmailSubject[string(constants.EmailTypeConfirmation)]
		body := constants.NotificationConstants.EmailBody[string(constants.EmailTypeConfirmation)]

		// Simulate sending email
		err := s.EmailService.SendEmail(subject, notification.Email, fmt.Sprintf(body, notification.Code))
		if err != nil {
			log.Printf("Error sending email to %s: %v", notification.Email, err)
			continue
		}

		// Update notification status to success
		err = s.NotificationRepo.UpdateStatus(notification.ID, "success")
		if err != nil {
			log.Printf("Error updating notification status for ID %s: %v", notification.ID, err)
			continue
		}

		log.Printf("Successfully sent email to %s and updated notification status", notification.Email)
	}

	doneChan <- true
}
