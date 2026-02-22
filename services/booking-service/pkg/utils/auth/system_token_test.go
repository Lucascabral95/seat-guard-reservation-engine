package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestGenerateSystemToken(t *testing.T) {
	t.Setenv("JWT_SECRET", "secret")
	tok, err := GenerateSystemToken()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	parsed, err := jwt.Parse(tok, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil || !parsed.Valid {
		t.Fatalf("expected valid token, err=%v", err)
	}

	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatalf("expected map claims")
	}
	if claims["sub"] != "sqs-worker-service" || claims["role"] != "system" {
		t.Fatalf("unexpected claims: %#v", claims)
	}
	exp, ok := claims["exp"].(float64)
	if !ok || int64(exp) <= time.Now().Unix() {
		t.Fatalf("expected future exp claim, got %#v", claims["exp"])
	}
}
