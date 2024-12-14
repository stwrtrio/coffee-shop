package handlers

import (
	"fmt"
	"net/http"

	"github.com/stwrtrio/coffee-shop/internal/domain/services"
	"github.com/stwrtrio/coffee-shop/models"
	"github.com/stwrtrio/coffee-shop/pkg/constants"
	"github.com/stwrtrio/coffee-shop/pkg/middlewares"
	"github.com/stwrtrio/coffee-shop/pkg/utils"

	"github.com/labstack/echo/v4"
)

type OrderHandler struct {
	OrderService services.OrderService
}

func NewOrderHandler(service services.OrderService) *OrderHandler {
	return &OrderHandler{service}
}

// Create Order
func (h *OrderHandler) CreateOrder(c echo.Context) error {
	var req *models.RequestOrder
	if err := c.Bind(&req); err != nil {
		fmt.Println("err1:", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	if err := c.Validate(req); err != nil {
		fmt.Println("err2:", err)
		return utils.FailResponse(c, http.StatusBadRequest, "Invalid request body")
	}

	claims, ok := middlewares.GetUserFromContext(c)
	if !ok {
		return utils.FailResponse(c, http.StatusBadRequest, constants.ErrInvalidToken)
	}

	staffID := claims.UserID

	if err := h.OrderService.CreateOrder(c.Request().Context(), req, staffID); err != nil {
		if err.Error() == constants.ErrorMenuIDNotFound {
			return utils.FailResponse(c, http.StatusBadRequest, err.Error())
		}
		return utils.FailResponse(c, http.StatusInternalServerError, "Failed to create order")
	}

	return utils.SuccessResponse(c, http.StatusCreated, "order created successfully", nil)
}
