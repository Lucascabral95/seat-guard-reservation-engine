package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestStripeHandler_CreateCartCheckoutSession_Validations(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("missing stripe key", func(t *testing.T) {
		t.Setenv("STRIPE_SECRET_KEY", "")
		r := gin.New()
		r.POST("/stripe", CreateCartCheckoutSession(nil, nil))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/stripe", bytes.NewBufferString(`{}`)))
		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w.Code)
		}
	})

	t.Run("bad body", func(t *testing.T) {
		t.Setenv("STRIPE_SECRET_KEY", "sk_test_x")
		r := gin.New()
		r.POST("/stripe", CreateCartCheckoutSession(nil, nil))
		req := httptest.NewRequest(http.MethodPost, "/stripe", bytes.NewBufferString("{"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", w.Code)
		}
	})

	t.Run("empty items", func(t *testing.T) {
		t.Setenv("STRIPE_SECRET_KEY", "sk_test_x")
		r := gin.New()
		r.POST("/stripe", CreateCartCheckoutSession(nil, nil))
		req := httptest.NewRequest(http.MethodPost, "/stripe", bytes.NewBufferString(`{"userId":"u1","currency":"usd","items":[]}`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", w.Code)
		}
	})
}

func TestBuildCheckoutResponse(t *testing.T) {
	resp := BuildCheckoutResponse(&ResponseCartCheckoutReq{OrderBookingId: "o1", UserId: "u1", Currency: "usd"})
	if resp["orderBookingId"] != "o1" || resp["userId"] != "u1" {
		t.Fatalf("unexpected response payload: %#v", resp)
	}
}
