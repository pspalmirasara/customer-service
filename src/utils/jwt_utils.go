package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))
var jwtIssuer = (os.Getenv("JWT_ISSUER"))

type CustomClaims struct {
	CustomerId string `json:"customerId"`
	jwt.RegisteredClaims
}

// GenerateJWT generates a JWT token with a given payload
func GenerateJWT(customerID interface{}) (string, error) {
	var claims CustomClaims

	if customerID == nil {
		claims = CustomClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    jwtIssuer,
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			},
		}
	} else {
		claims = CustomClaims{
			CustomerId: fmt.Sprintf("%v", customerID),
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    jwtIssuer,
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			},
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtSecret)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
