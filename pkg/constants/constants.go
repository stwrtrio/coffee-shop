package constants

type EmailType string

const (
	// Handler Error Message
	ErrInvalidRequestBody = "Invalid request body"
	ErrInvalidToken       = "Invalid Token"

	// Email Error Message
	EmailTypeConfirmation  EmailType = "confirmation"
	EmailTypePasswordReset EmailType = "password_reset"

	// Menu Error Message
	ErrMenuAlreadyExists = "menu already exist"
	ErrCategoryNotExists = "category is not exist"

	// Redis Cache Key
	MenusCacheKey = "menus:page=%d:limit=%d"
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
