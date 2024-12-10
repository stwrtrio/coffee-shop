package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/stwrtrio/coffee-shop/internal/domain/services"
	"github.com/stwrtrio/coffee-shop/models"
	"github.com/stwrtrio/coffee-shop/pkg/constants"
	"github.com/stwrtrio/coffee-shop/pkg/middlewares"
	"github.com/stwrtrio/coffee-shop/pkg/utils"
)

type MenuHandler struct {
	service services.MenuService
}

func NewMenuHandler(service services.MenuService) *MenuHandler {
	return &MenuHandler{service: service}
}

func (h *MenuHandler) CreateMenu(c echo.Context) error {
	ctx := c.Request().Context()
	var req *models.MenuRequest
	if err := c.Bind(&req); err != nil {
		return utils.FailResponse(c, http.StatusBadRequest, constants.ErrInvalidRequestBody)
	}

	claims, ok := middlewares.GetUserFromContext(c)
	if !ok {
		return utils.FailResponse(c, http.StatusBadRequest, constants.ErrInvalidToken)
	}

	req.CreatedBy = claims.UserID

	menu, err := h.service.CreateMenu(ctx, req)
	if err != nil {
		if err.Error() == constants.ErrCategoryNotExists || err.Error() == constants.ErrMenuAlreadyExists {
			return utils.FailResponse(c, http.StatusBadRequest, err.Error())
		}
		return utils.FailResponse(c, http.StatusInternalServerError, "Failed to create menu")
	}

	return utils.SuccessResponse(c, http.StatusCreated, "", menu)
}
