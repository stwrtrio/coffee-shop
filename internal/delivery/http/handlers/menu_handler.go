package handlers

import (
	"net/http"
	"strconv"

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

func (h *MenuHandler) GetAllMenus(c echo.Context) error {
	useCache := c.QueryParam("use-cache")
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil || limit < 1 {
		limit = 10
	}

	menus, err := h.service.GetAllMenus(c.Request().Context(), page, limit, useCache)
	if err != nil {
		return utils.FailResponse(c, http.StatusInternalServerError, "Failed to fetch menus")
	}

	return utils.SuccessResponse(c, http.StatusOK, "", menus)
}

func (h *MenuHandler) UpdateMenu(c echo.Context) error {
	var req *models.MenuRequest
	if err := c.Bind(&req); err != nil {
		return utils.FailResponse(c, http.StatusBadRequest, constants.ErrInvalidRequestBody)
	}

	req.MenuID = c.Param("id")

	claims, ok := middlewares.GetUserFromContext(c)
	if !ok {
		return utils.FailResponse(c, http.StatusBadRequest, constants.ErrInvalidToken)
	}

	req.UpdatedBy = claims.UserID

	menu, err := h.service.UpdateMenu(c.Request().Context(), req)
	if err != nil {
		return utils.FailResponse(c, http.StatusInternalServerError, "Failed to update menu")
	}

	return utils.SuccessResponse(c, http.StatusOK, "", menu)
}
