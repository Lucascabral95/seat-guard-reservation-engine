package handlers

import (
	"booking-service/internal/models"
	"booking-service/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type EventHandler struct {
	service *services.EventService
}

// Constructor
func NewEventHandler(service *services.EventService) *EventHandler {
	return &EventHandler{service: service}
}

// POST /events
func (h *EventHandler) CreateEvent(c *gin.Context) {
	var event models.Event

	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format: " + err.Error()})
		return
	}

	if err := h.service.CreateEvent(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, event)
}

func (h *EventHandler) GetEventByID(c *gin.Context) {
	id := c.Param("id")

	_, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	event, err := h.service.GetEvent(id)
	if err != nil {
		if err.Error() == "event not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, event)
}

// GET /events
func (h *EventHandler) GetAllEvents(c *gin.Context) {
	events, err := h.service.GetAllEvents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch events"})
		return
	}

	c.JSON(http.StatusOK, events)
}

// PUT /events/:id
func (h *EventHandler) UpdateEvent(c *gin.Context) {
	id := c.Param("id")
	var event models.Event

	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateEvent(id, &event); err != nil {
		if err.Error() == "cannot update: event not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event updated successfully"})
}

// DELETE /events/:id
func (h *EventHandler) DeleteEvent(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.DeleteEvent(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete event"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event deleted successfully"})
}

// PATCH /events/availability/:id
func (h *EventHandler) UpdateAvailabilityForEvent(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.UpdateEventAvailability(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update availability for event"})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event availability successfully updated"})
}
