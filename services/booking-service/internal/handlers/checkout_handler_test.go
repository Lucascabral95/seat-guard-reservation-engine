package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCheckoutHandler_Create_ValidationErrors(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &CheckoutHandler{}
	r := gin.New()
	r.POST("/checkouts", h.Create)

	t.Run("bad json", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/checkouts", bytes.NewBufferString("{"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", w.Code)
		}
	})

	t.Run("missing fields", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/checkouts", bytes.NewBufferString(`{"orderId":""}`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", w.Code)
		}
	})
}
