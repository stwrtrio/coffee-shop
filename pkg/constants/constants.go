package constants

type EmailType string

const (
	EmailTypeConfirmation  EmailType = "confirmation"
	EmailTypePasswordReset EmailType = "password_reset"
)

var NotificationConstants = struct {
	EmailType    map[string]EmailType
	EmailSubject map[string]string
	EmailBody    map[string]string
}{
	EmailType: map[string]EmailType{
		"confirmation":   EmailTypeConfirmation,
		"password_reset": EmailTypePasswordReset,
	},
	EmailSubject: map[string]string{
		"confirmation":   "Your Confirmation Code",
		"password_reset": "Password Reset Request",
	},
	EmailBody: map[string]string{
		"confirmation":   "Your confirmation code is: %s",
		"password_reset": "Click here to reset your password: %s",
	},
}
