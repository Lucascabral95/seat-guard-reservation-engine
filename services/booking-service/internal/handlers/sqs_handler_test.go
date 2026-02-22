package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestSQSHandler_Send_ValidationPaths(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &SQSHandler{}
	r := gin.New()
	r.POST("/sqs", h.Send)

	t.Run("bad json", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/sqs", bytes.NewBufferString("{"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", w.Code)
		}
	})

	t.Run("missing metadata ignored", func(t *testing.T) {
		payload := `{"id":"evt_1","type":"checkout.session.completed","data":{"object":{"payment_status":"paid","metadata":{"user_id":"","seat_ids":"","order_id":""}}}}`
		req := httptest.NewRequest(http.MethodPost, "/sqs", bytes.NewBufferString(payload))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
	})
}
