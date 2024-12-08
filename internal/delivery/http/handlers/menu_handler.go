package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/stwrtrio/coffee-shop/internal/domain/services"
)

type MenuHandler struct {
	service *services.MenuService
}

func NewMenuHandler(service *services.MenuService) *MenuHandler {
	return &MenuHandler{service: service}
}

func (h *MenuHandler) CreateMenu(c echo.Context) error {
	return nil
}
