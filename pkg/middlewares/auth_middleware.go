package middlewares

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/stwrtrio/coffee-shop/pkg/helpers"
	"github.com/stwrtrio/coffee-shop/pkg/utils"
)

func GetUserFromContext(c echo.Context) (*helpers.Claims, bool) {
	claims, ok := c.Get("user").(*helpers.Claims)
	return claims, ok
}

func JWTMiddleware(config utils.JwtConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return utils.FailResponse(c, http.StatusUnauthorized, "Missing Authorization header")
			}

			tokenString := authHeader[len("Bearer "):]

			// Validate the token
			claims, err := helpers.ValidateJWTToken(&config, tokenString)
			if err != nil {
				return utils.FailResponse(c, http.StatusUnauthorized, "invalid or expired token")
			}

			c.Set("user", claims)
			c.Set("user_role", claims.Role)
			return next(c)
		}
	}
}

func RoleMiddleware(allowedRoles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Retrieve the user's role from the context (assumes JWT middleware sets this)
			role := c.Get("user_role")
			if role == nil {
				return utils.FailResponse(c, http.StatusUnauthorized, "Unauthorized")
			}

			// Check if the role is allowed
			for _, allowedRole := range allowedRoles {
				if role == allowedRole {
					return next(c)
				}
			}

			return utils.FailResponse(c, http.StatusForbidden, "Access Forbiden")
		}
	}
}
