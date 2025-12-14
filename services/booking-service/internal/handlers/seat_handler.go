package handlers

import (
	"booking-service/internal/models"
	"booking-service/internal/services"
	"booking-service/pkg/utils"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SeatHandler struct {
	service *services.SeatService
}

func NewSeatHandler(service *services.SeatService) *SeatHandler {
	return &SeatHandler{service: service}
}

// POST /seats
func (h *SeatHandler) CreateSeat(c *gin.Context) {
	var seat models.Seat

	if err := c.ShouldBindJSON(&seat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format: " + err.Error()})
		return
	}

	if err := h.service.CreateSeat(&seat); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create seat"})
		return
	}

	c.JSON(http.StatusCreated, seat)
}

// GET /seats
func (h *SeatHandler) GetSeats(c *gin.Context) {
	seats, err := h.service.GetSeats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch seats"})
		return
	}

	c.JSON(http.StatusOK, seats)
}

// GET /seats/:id
func (h *SeatHandler) GetSeat(c *gin.Context) {
	id := c.Param("id")

	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	seat, err := h.service.GetSeat(id)
	if err != nil {
		if errors.Is(err, utils.ErrSeatNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Seat not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch seat"})
		}
		return
	}

	c.JSON(http.StatusOK, seat)
}

// GET /seats/event/:eventId
func (h *SeatHandler) GetSeatsByEventId(c *gin.Context) {
	eventId := c.Param("eventId")

	_, err := uuid.Parse(eventId)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Invalid UUID format"})
		return
	}

	seat, err := h.service.GetSeatByEventId(eventId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch seats by event id"})
		return
	}

	c.JSON(http.StatusOK, seat)
}

// PATCH /seats/:id
func (h *SeatHandler) UpdateSeat(c *gin.Context) {
	id := c.Param("id")
	var seat models.Seat

	if err := c.ShouldBindJSON(&seat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format: " + err.Error()})
		return
	}

	if err := h.service.UpdateSeatStatus(id, seat.Status); err != nil {
		if errors.Is(err, utils.ErrSeatNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Seat not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update seat"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Seat updated successfully"})
}

// Bloquear por 15 minutos. Requisitos: recibir uid y id del asiento
func (h *SeatHandler) LockSeat(c *gin.Context) {
	id := c.Param("id")
	uid := c.Param("uid")

	if err := h.service.LockSeat(id, uid); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Seat locked successfully",
		"expiresAt": time.Now().Add(15 * time.Minute).Format(time.RFC3339),
	})
}
