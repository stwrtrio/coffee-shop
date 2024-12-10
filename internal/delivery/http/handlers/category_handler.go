package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/stwrtrio/coffee-shop/internal/domain/services"
	"github.com/stwrtrio/coffee-shop/models"
	"github.com/stwrtrio/coffee-shop/pkg/utils"
)

type CategoryHandler struct {
	CategoryService services.CategoryService
}

func NewCategoryHandler(service services.CategoryService) *CategoryHandler {
	return &CategoryHandler{service}
}

func (h *CategoryHandler) CreateCategory(c echo.Context) error {
	ctx := c.Request().Context()
	var req *models.CategoryRequest
	if err := c.Bind(&req); err != nil {
		return utils.FailResponse(c, http.StatusBadRequest, "Invalid request body")
	}

	result, err := h.CategoryService.CreateCategory(ctx, req)
	if err != nil {
		return utils.FailResponse(c, http.StatusInternalServerError, "Failed to create category")
	}

	return utils.SuccessResponse(c, http.StatusCreated, "category create successfully", result)
}
