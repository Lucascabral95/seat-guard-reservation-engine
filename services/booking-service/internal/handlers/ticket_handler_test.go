package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestTicketHandler_GetTicketMetadata_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &TicketHandler{}
	r := gin.New()
	r.GET("/tickets/:orderID", h.GetTicketMetadata)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/tickets/o1", nil))
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestTicketHandler_CreateTicketFromEndpoint_BadJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &TicketHandler{}
	r := gin.New()
	r.POST("/tickets", h.CreateTicketFromEndpoint)

	req := httptest.NewRequest(http.MethodPost, "/tickets", bytes.NewBufferString("{"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}
