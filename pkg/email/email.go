package email

import (
	"fmt"
	"net/smtp"
)

// EmailService defines the interface for sending emails.
type EmailService interface {
	SendEmail(subject string, to string, code string) error
}

// emailService is the implementation of EmailService.
type emailService struct {
	SMTPHost     string
	SMTPPort     string
	SenderEmail  string
	SenderPasswd string
}

// NewEmailService initializes a new EmailService.
func NewEmailService(smtpHost, smtpPort, senderEmail, senderPasswd string) EmailService {
	return &emailService{
		SMTPHost:     smtpHost,
		SMTPPort:     smtpPort,
		SenderEmail:  senderEmail,
		SenderPasswd: senderPasswd,
	}
}

// SendEmail sends an email to the specified recipient with a code.
func (e *emailService) SendEmail(subject string, to string, body string) error {
	// Define the email body
	emailSubject := fmt.Sprintf("Subject: %s\n", subject)
	message := emailSubject + "\n" + body

	// Set up authentication information
	auth := smtp.PlainAuth("", e.SenderEmail, e.SenderPasswd, e.SMTPHost)

	// Send the email
	err := smtp.SendMail(
		fmt.Sprintf("%s:%s", e.SMTPHost, e.SMTPPort), // Address of the SMTP server
		auth,            // Auth information
		e.SenderEmail,   // From email
		[]string{to},    // To email
		[]byte(message), // Email message
	)

	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}
