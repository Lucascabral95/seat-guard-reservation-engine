package handlers

import (
	"booking-service/pkg/domain"
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type mockEmailService struct {
	sendAsyncFn        func(*domain.Email) error
	sendBulkFn         func([]*domain.Email)
	sendPurchaseEmailFn func(context.Context, string, string, string, float64) error
}

func (m *mockEmailService) SendAsync(e *domain.Email) error { return m.sendAsyncFn(e) }
func (m *mockEmailService) SendBulk(e []*domain.Email)      { m.sendBulkFn(e) }
func (m *mockEmailService) Shutdown()                        {}
func (m *mockEmailService) SendPurchaseEmail(ctx context.Context, to, name, orderId string, amount float64) error {
	return m.sendPurchaseEmailFn(ctx, to, name, orderId, amount)
}

func TestEmailHandler_SendAsync_QueueFull(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewEmailHandler(&mockEmailService{sendAsyncFn: func(*domain.Email) error { return errors.New("full") }})
	r := gin.New()
	r.POST("/emails/send-bulk-async", h.SendAsync)

	req := httptest.NewRequest(http.MethodPost, "/emails/send-bulk-async", bytes.NewBufferString(`{"to":["a@a.com"],"subject":"s","body":"b"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503, got %d", w.Code)
	}
}

func TestEmailHandler_SendBulk_Accepted(t *testing.T) {
	gin.SetMode(gin.TestMode)
	called := false
	h := NewEmailHandler(&mockEmailService{sendBulkFn: func(emails []*domain.Email) { called = len(emails) == 1 }})
	r := gin.New()
	r.POST("/emails/send-bulk", h.SendBulk)

	req := httptest.NewRequest(http.MethodPost, "/emails/send-bulk", bytes.NewBufferString(`{"emails":[{"to":["a@a.com"],"subject":"s","body":"b"}]}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusAccepted || !called {
		t.Fatalf("expected 202 and bulk call, got code=%d called=%v", w.Code, called)
	}
}
