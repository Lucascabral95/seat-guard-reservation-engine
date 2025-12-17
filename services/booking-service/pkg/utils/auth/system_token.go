package auth

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateSystemToken crea un token firmado para uso interno entre microservicios
func GenerateSystemToken() (string, error) {
	claims := jwt.MapClaims{
		"sub":  "sqs-worker-service",
		"role": "system",                               // ESTO es lo que busca el middleware
		"exp":  time.Now().Add(time.Minute * 5).Unix(), // Token de vida corta (5 min)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}
