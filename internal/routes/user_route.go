package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/stwrtrio/coffee-shop/internal/delivery/http/handlers"
)

// RegisterAuthRoutes sets up routes for authentication-related endpoints.
func RegisterAuthRoutes(e *echo.Echo, customerHandler *handlers.CustomerHandler) {
	// Public routes
	e.POST("/api/customer/register", customerHandler.RegisterCustomer)
	e.POST("/api/customer/login", customerHandler.LoginCustomer)
}
