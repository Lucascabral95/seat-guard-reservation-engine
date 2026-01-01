package handlers

import (
	"booking-service/internal/services"
	"booking-service/pkg/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type EmailHandler struct {
	service services.EmailService
}

func NewEmailHandler(service services.EmailService) *EmailHandler {
	return &EmailHandler{service: service}
}

type SendRequest struct {
	To      []string `json:"to" binding:"required"`
	Subject string   `json:"subject" binding:"required"`
	Body    string   `json:"body" binding:"required"`
}

type BulkRequest struct {
	Emails []SendRequest `json:"emails" binding:"required,min=1"`
}

// POST /send-sync - Envío síncrono
type SendPurchaseRequest struct {
	To      string  `json:"to" binding:"required,email"`
	Name    string  `json:"name" binding:"required"`
	OrderId string  `json:"orderId" binding:"required"`
	Amount  float64 `json:"amount" binding:"required"`
}

func (h *EmailHandler) SendSync(c *gin.Context) {
	var req SendPurchaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	if err := h.service.SendPurchaseEmail(ctx, req.To, req.Name, req.OrderId, req.Amount); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "sent"})
}

// POST /send - Envío asíncrono
func (h *EmailHandler) SendAsync(c *gin.Context) {
	var req SendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	email := &domain.Email{
		To:      req.To,
		Subject: req.Subject,
		Body:    req.Body,
	}

	if err := h.service.SendAsync(email); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "queue full"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"status": "queued"})
}

// POST /send-bulk - Envío masivo
func (h *EmailHandler) SendBulk(c *gin.Context) {
	var req BulkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var emails []*domain.Email
	for _, item := range req.Emails {
		emails = append(emails, &domain.Email{
			To:      item.To,
			Subject: item.Subject,
			Body:    item.Body,
		})
	}

	h.service.SendBulk(emails)

	c.JSON(http.StatusAccepted, gin.H{
		"status": "processing",
		"total":  len(emails),
	})
}
