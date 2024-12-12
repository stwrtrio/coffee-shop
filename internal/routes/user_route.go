package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/stwrtrio/coffee-shop/internal/delivery/http/handlers"
	"github.com/stwrtrio/coffee-shop/pkg/middlewares"
	"github.com/stwrtrio/coffee-shop/pkg/utils"
)

// RegisterUserRoutes sets up routes for authentication-related endpoints.
func RegisterUserRoutes(e *echo.Echo,
	config *utils.Config,
	customerHandler *handlers.UserHandler) {
	// Public routes
	userGroup := e.Group("/api/user")
	userGroup.POST("/register", customerHandler.RegisterUser)
	userGroup.POST("/login", customerHandler.LoginUser)
	userGroup.POST("/confirm-code", customerHandler.ConfirmCode)

	userGroup.Use(middlewares.JWTMiddleware(config.Jwt))
	userGroup.PUT("/update/:id", customerHandler.UpdateUser)
}
