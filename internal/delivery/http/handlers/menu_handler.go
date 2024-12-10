package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/stwrtrio/coffee-shop/internal/domain/services"
	"github.com/stwrtrio/coffee-shop/models"
	"github.com/stwrtrio/coffee-shop/pkg/utils"
)

type MenuHandler struct {
	service services.MenuService
}

func NewMenuHandler(service services.MenuService) *MenuHandler {
	return &MenuHandler{service: service}
}

func (h *MenuHandler) CreateMenu(c echo.Context) error {
	var menu models.Menu
	if err := c.Bind(&menu); err != nil {
		return utils.FailResponse(c, http.StatusBadRequest, "Invalid request body")
	}

	if err := h.service.CreateMenu(&menu); err != nil {
		return utils.FailResponse(c, http.StatusInternalServerError, "Failed to create menu")
	}

	return utils.SuccessResponse(c, http.StatusCreated, "menu created successfully", nil)
}
