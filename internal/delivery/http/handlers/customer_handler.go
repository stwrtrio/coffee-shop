package handlers

import (
	"net/http"

	"github.com/stwrtrio/coffee-shop/internal/domain/services"
	"github.com/stwrtrio/coffee-shop/models"
	"github.com/stwrtrio/coffee-shop/pkg/utils"
	"gorm.io/gorm"

	"github.com/labstack/echo/v4"
)

type CustomerHandler struct {
	CustomerService services.CustomerService
}

func NewCustomerHandler(service services.CustomerService) *CustomerHandler {
	return &CustomerHandler{service}
}

// Register a new customer
func (h *CustomerHandler) RegisterCustomer(c echo.Context) error {
	ctx := c.Request().Context()
	var req *models.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return utils.FailResponse(c, http.StatusBadRequest, "Invalid request body")
	}

	if err := h.CustomerService.RegisterCustomer(ctx, req); err != nil {
		return utils.FailResponse(c, http.StatusInternalServerError, "Failed to create customer")
	}

	return utils.SuccessResponse(c, http.StatusCreated, "customer registered successfully", nil)
}

// Login handles customer login
func (h *CustomerHandler) LoginCustomer(c echo.Context) error {
	ctx := c.Request().Context()

	// Parse request body
	var customerRequest models.LoginRequest
	if err := c.Bind(&customerRequest); err != nil {
		return utils.FailResponse(c, http.StatusBadRequest, "Invalid request body")
	}

	// Call the service to validate credentials and generate token
	token, err := h.CustomerService.LoginCustomer(ctx, &customerRequest)
	if err != nil {
		if err.Error() == "customer is not confirmed" {
			return utils.FailResponse(c, http.StatusUnauthorized, "Customer email has not been confirmed.")
		}
		return utils.FailResponse(c, http.StatusUnauthorized, "invalid credentials")
	}

	return utils.SuccessResponse(c, http.StatusOK, "Access granted", map[string]string{"token": token})
}

// ConfirmCode verifies the customer's confirmation code
func (h *CustomerHandler) ConfirmCode(c echo.Context) error {
	ctx := c.Request().Context()
	var req models.ConfirmCodeRequest
	if err := c.Bind(&req); err != nil {
		return utils.FailResponse(c, http.StatusBadRequest, "Invalid request body")
	}

	if err := c.Validate(req); err != nil {
		return utils.FailResponse(c, http.StatusBadRequest, "Invalid request body")
	}

	err := h.CustomerService.ConfirmCode(ctx, req.Email, req.Code)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.FailResponse(c, http.StatusBadRequest, "customer not found")
		}
		if err.Error() == "invalid confirmation code" || err.Error() == "confirmation code expired" {
			return utils.FailResponse(c, http.StatusBadRequest, err.Error())
		}
		return utils.FailResponse(c, http.StatusInternalServerError, "Confirmation failed")
	}

	return utils.SuccessResponse(c, http.StatusOK, "Email confirmed successfully", nil)
}
