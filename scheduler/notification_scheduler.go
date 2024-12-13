package scheduler

import (
	"fmt"
	"log"

	"github.com/stwrtrio/coffee-shop/models"
	"github.com/stwrtrio/coffee-shop/pkg/constants"
)

func (s *scheduler) processNotifications(notificationChan <-chan models.Notification, doneChan chan<- bool) {
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
