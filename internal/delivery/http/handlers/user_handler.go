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

type UserHandler struct {
	UserService services.UserService
}

func NewUserHandler(service services.UserService) *UserHandler {
	return &UserHandler{service}
}

// Register a new user
func (h *UserHandler) RegisterUser(c echo.Context) error {
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

	// Create user
	user := &models.User{
		ID:           uuid.NewString(),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
	}

	if err := h.UserService.RegisterUser(ctx, user); err != nil {
		return utils.FailResponse(c, http.StatusInternalServerError, "Failed to create user")
	}

	return utils.SuccessResponse(c, http.StatusCreated, "user registered successfully", nil)
}
