package jwt

import (
	"errors"
	"fmt"
	"os"
	"product-api/internal/config"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Claims represents the JWT claims
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateToken generates a JWT token with a given username.
func GenerateToken() (string, error) {
	// Read JWT secret from environment variables
	envs := config.LoadConfig()

	secretKey := envs.TOKEN.Secret

	expirationTime := envs.TOKEN.Expiry

	fmt.Println(expirationTime)
	// Parse expiration time
	expirationDuration, err := time.ParseDuration(expirationTime)
	if err != nil {
		return "", errors.New("invalid JWT expiration time format")
	}

	username := envs.TOKEN.User

	// Create claims
	claims := Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "product-api",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expirationDuration)),
		},
	}

	// Create the token with HS256 signing method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token and return it
	return token.SignedString([]byte(secretKey))
}

// ValidateToken validates a JWT token and returns the claims.
func ValidateToken(tokenStr string) (*Claims, error) {
	// Read the secret key from env
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		return nil, errors.New("missing JWT secret key in environment")
	}

	// Parse the token and validate it
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(secretKey), nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid or expired token")
	}

	// Return the claims if the token is valid
	if claims, ok := token.Claims.(*Claims); ok {
		return claims, nil
	}
	return nil, errors.New("could not parse token claims")
}
