package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/stwrtrio/coffee-shop/internal/delivery/http/handlers"
)

// RegisterAuthRoutes sets up routes for authentication-related endpoints.
func RegisterAuthRoutes(e *echo.Echo, userHandler *handlers.UserHandler) {
	// Public routes
	e.POST("/api/user/register", userHandler.RegisterUser)
}
