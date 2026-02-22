package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestSeatHandler_GetSeat_InvalidUUID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &SeatHandler{}
	r := gin.New()
	r.GET("/seats/:id", h.GetSeat)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/seats/bad", nil))
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestSeatHandler_GetSeatsByEventId_InvalidUUID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &SeatHandler{}
	r := gin.New()
	r.GET("/seats/event/:eventId", h.GetSeatsByEventId)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/seats/event/bad", nil))
	if w.Code != http.StatusBadGateway {
		t.Fatalf("expected 502, got %d", w.Code)
	}
}
