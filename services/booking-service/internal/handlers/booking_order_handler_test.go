package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestBookingOrderHandler_CreateBookingOrder_BadJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &BookingOrderHandler{}
	r := gin.New()
	r.POST("/booking-orders", h.CreateBookingOrder)

	req := httptest.NewRequest(http.MethodPost, "/booking-orders", bytes.NewBufferString("{"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestBookingOrderHandler_GetBookingOrderById_InvalidUUID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &BookingOrderHandler{}
	r := gin.New()
	r.GET("/booking-orders/:id", h.GetBookingOrderById)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/booking-orders/not-uuid", nil))
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestBookingOrderHandler_UpdateBookingOrder_InvalidStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &BookingOrderHandler{}
	r := gin.New()
	r.PATCH("/booking-orders/:id", h.UpdateBookingOrder)

	body := `{"status":"BAD_STATUS"}`
	req := httptest.NewRequest(http.MethodPatch, "/booking-orders/11111111-1111-1111-1111-111111111111", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}
