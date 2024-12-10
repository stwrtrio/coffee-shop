package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/stwrtrio/coffee-shop/internal/delivery/http/handlers"
)

// RegisterUserRoutes sets up routes for authentication-related endpoints.
func RegisterUserRoutes(e *echo.Echo, customerHandler *handlers.UserHandler) {
	// Public routes
	e.POST("/api/user/register", customerHandler.RegisterUser)
	e.POST("/api/user/login", customerHandler.LoginUser)
	e.POST("/api/user/confirm-code", customerHandler.ConfirmCode)
}
