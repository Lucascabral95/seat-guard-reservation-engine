package utils

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func TestExtractToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = httptest.NewRequest(http.MethodGet, "/", nil)

	_, err := ExtractToken(ctx)
	if err == nil {
		t.Fatalf("expected missing header error")
	}

	ctx.Request.Header.Set("Authorization", "Token abc")
	_, err = ExtractToken(ctx)
	if !errors.Is(err, jwt.ErrTokenMalformed) {
		t.Fatalf("expected malformed token error, got %v", err)
	}

	ctx.Request.Header.Set("Authorization", "Bearer abc.def.ghi")
	tok, err := ExtractToken(ctx)
	if err != nil || tok != "abc.def.ghi" {
		t.Fatalf("expected token extraction, got token=%q err=%v", tok, err)
	}
}
