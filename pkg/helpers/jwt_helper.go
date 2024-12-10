package helpers

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/stwrtrio/coffee-shop/pkg/utils"
)

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	JTI    string `json:"jti"`
	jwt.RegisteredClaims
}

// GenerateJWTToken creates a JWT token for the user
func GenerateJWTToken(config *utils.JwtConfig, userID, email, role string) (string, error) {
	expirationTime := time.Now().Add(config.Expiry)
	claims := Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		JTI: uuid.NewString(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.SecretKey))
}

// ValidateJWTToken parses and validates a JWT token
func ValidateJWTToken(config *utils.JwtConfig, tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.SecretKey), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, jwt.ErrInvalidKey
	}

	return claims, nil
}
