package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/stwrtrio/coffee-shop/internal/delivery/http/handlers"
	"github.com/stwrtrio/coffee-shop/pkg/middlewares"
	"github.com/stwrtrio/coffee-shop/pkg/utils"
)

func RegisterOrderRoutes(e *echo.Echo,
	config *utils.Config,
	orderHandler *handlers.OrderHandler) {

	staffGroup := e.Group("/api/staff")
	staffGroup.Use(middlewares.JWTMiddleware(config.Jwt))
	staffGroup.Use(middlewares.RoleMiddleware(config.RolesAllowed...))

	// Order routes
	staffGroup.POST("/order", orderHandler.CreateOrder)
}
