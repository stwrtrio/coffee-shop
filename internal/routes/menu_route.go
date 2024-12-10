package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/stwrtrio/coffee-shop/internal/delivery/http/handlers"
	"github.com/stwrtrio/coffee-shop/pkg/middlewares"
	"github.com/stwrtrio/coffee-shop/pkg/utils"
)

// RegisterMenuRoutes sets up routes for authentication-related endpoints.
func RegisterMenuRoutes(e *echo.Echo,
	config *utils.Config,
	menuHandler *handlers.MenuHandler,
	categoryHandler *handlers.CategoryHandler) {

	staffGroup := e.Group("/api/staff")
	staffGroup.Use(middlewares.JWTMiddleware(config.Jwt))
	staffGroup.Use(middlewares.RoleMiddleware(config.RolesAllowed...))

	// Menu routes
	staffGroup.POST("/menu", menuHandler.CreateMenu)

	// Category routes
	staffGroup.POST("/category", categoryHandler.CreateCategory)
}
