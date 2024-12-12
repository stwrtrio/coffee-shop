package handlers

import (
	"net/http"

	"github.com/stwrtrio/coffee-shop/internal/domain/services"
	"github.com/stwrtrio/coffee-shop/models"
	"github.com/stwrtrio/coffee-shop/pkg/utils"
	"gorm.io/gorm"

	"github.com/labstack/echo/v4"
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
	var req *models.UserRegisterRequest
	if err := c.Bind(&req); err != nil {
		return utils.FailResponse(c, http.StatusBadRequest, "Invalid request body")
	}

	if err := h.UserService.RegisterUser(ctx, req); err != nil {
		return utils.FailResponse(c, http.StatusInternalServerError, "Failed to create user")
	}

	return utils.SuccessResponse(c, http.StatusCreated, "user registered successfully", nil)
}

// Login handles user login
func (h *UserHandler) LoginUser(c echo.Context) error {
	ctx := c.Request().Context()

	// Parse request body
	var userRequest models.UserLoginRequest
	if err := c.Bind(&userRequest); err != nil {
		return utils.FailResponse(c, http.StatusBadRequest, "Invalid request body")
	}

	// Call the service to validate credentials and generate token
	token, err := h.UserService.LoginUser(ctx, &userRequest)
	if err != nil {
		if err.Error() == "user is not confirmed" {
			return utils.FailResponse(c, http.StatusUnauthorized, "User email has not been confirmed.")
		}
		return utils.FailResponse(c, http.StatusUnauthorized, "invalid credentials")
	}

	return utils.SuccessResponse(c, http.StatusOK, "Access granted", map[string]string{"token": token})
}

// ConfirmCode verifies the user's confirmation code
func (h *UserHandler) ConfirmCode(c echo.Context) error {
	ctx := c.Request().Context()
	var req models.UserConfirmCodeRequest
	if err := c.Bind(&req); err != nil {
		return utils.FailResponse(c, http.StatusBadRequest, "Invalid request body")
	}

	if err := c.Validate(req); err != nil {
		return utils.FailResponse(c, http.StatusBadRequest, "Invalid request body")
	}

	err := h.UserService.ConfirmCode(ctx, req.Email, req.Code)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.FailResponse(c, http.StatusBadRequest, "user not found")
		}
		if err.Error() == "invalid confirmation code" || err.Error() == "confirmation code expired" {
			return utils.FailResponse(c, http.StatusBadRequest, err.Error())
		}
		return utils.FailResponse(c, http.StatusInternalServerError, "Confirmation failed")
	}

	return utils.SuccessResponse(c, http.StatusOK, "Email confirmed successfully", nil)
}

func (h *UserHandler) UpdateUser(c echo.Context) error {
	ctx := c.Request().Context()

	var req models.UserUpdateRequest
	if err := c.Bind(&req); err != nil {
		return utils.FailResponse(c, http.StatusBadRequest, "Invalid request body")
	}

	req.UserID = c.Param("id")

	user, err := h.UserService.UpdateUser(ctx, req)
	if err != nil {
		return utils.FailResponse(c, http.StatusInternalServerError, "Update User Failed")
	}

	return utils.SuccessResponse(c, http.StatusOK, "", user)
}
