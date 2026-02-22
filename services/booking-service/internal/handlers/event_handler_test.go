package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestEventHandler_GetEventByID_InvalidUUID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &EventHandler{}
	r := gin.New()
	r.GET("/events/:id", h.GetEventByID)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/events/bad", nil))
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestEventHandler_CreateEvent_BadJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &EventHandler{}
	r := gin.New()
	r.POST("/events", h.CreateEvent)

	req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewBufferString("{"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}
