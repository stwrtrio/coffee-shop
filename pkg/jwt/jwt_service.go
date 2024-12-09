package jwt

import (
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/stwrtrio/coffee-shop/pkg/utils"
)

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	JTI    string `json:"jti"`
	jwt.RegisteredClaims
}

// LoadJWTExpiry load JWT_EXPIRY in .env
func LoadJWTExpiry() time.Duration {
	expiryStr := os.Getenv("JWT_EXPIRY")
	expiry, err := time.ParseDuration(expiryStr)
	if err != nil {
		log.Fatalf("Invalid JWT_EXPIRY in .env: %v", err)
	}
	return expiry
}

// GenerateJWTToken creates a JWT token for the user
func GenerateJWTToken(config *utils.Config, userID, email string) (string, error) {
	expirationTime := time.Now().Add(config.Jwt.Expiry)
	claims := Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		JTI: uuid.NewString(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Jwt.SecretKey))
}

// ValidateJWTToken parses and validates a JWT token
func ValidateJWTToken(config *utils.Config, tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Jwt.SecretKey), nil
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
