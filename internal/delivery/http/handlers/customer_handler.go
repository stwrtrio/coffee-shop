package handlers

import (
	"net/http"

	"github.com/stwrtrio/coffee-shop/internal/domain/services"
	"github.com/stwrtrio/coffee-shop/models"
	"github.com/stwrtrio/coffee-shop/pkg/utils"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
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
	var req models.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return utils.FailResponse(c, http.StatusBadRequest, "Invalid request body")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return utils.FailResponse(c, http.StatusInternalServerError, "Failed to hash password")
	}

	// Create customer
	customer := &models.Customer{
		ID:                    uuid.NewString(),
		FirstName:             req.FirstName,
		LastName:              req.LastName,
		Email:                 req.Email,
		PasswordHash:          string(hashedPassword),
		EmailConfirmationCode: utils.GenerateConfirmationCode(),
		IsEmailConfirmed:      false,
	}

	if err := h.CustomerService.RegisterCustomer(ctx, customer); err != nil {
		return utils.FailResponse(c, http.StatusInternalServerError, "Failed to create customer")
	}

	return utils.SuccessResponse(c, http.StatusCreated, "customer registered successfully", nil)
}

// Login handles user login
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
		return utils.FailResponse(c, http.StatusUnauthorized, "invalid credentials")
	}

	return utils.SuccessResponse(c, http.StatusOK, "Access granted", map[string]string{"token": token})
}
