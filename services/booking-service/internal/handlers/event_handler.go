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

// CreateEvent godoc
// @Summary Crear evento
// @Description Crear un evento nuevo
// @Tags events
// @Accept json
// @Produce json
// @Param event body models.Event true "Datos del evento"
// @Success 201 {object} models.Event "Evento creado exitosamente"
// @Failure 400 {object} map[string]string "Formato JSON inválido"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /events [post]
// @Security BearerAuth
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

// GetEventByID godoc
// @Summary Obtener evento por ID
// @Description Obtener un evento por su ID
// @Tags events
// @Accept json
// @Produce json
// @Param id path string true "ID del evento"
// @Success 200 {object} models.Event
// @Failure 400 {object} map[string]string "Formato UUID inválido"
// @Failure 401 {object} map[string]string "No autorizado"
// @Failure 404 {object} map[string]string "Evento no encontrado"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /events/{id} [get]
// @Security BearerAuth
// GET /events/id
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

// GetAllEvents godoc
// @Summary Obtener todos los eventos
// @Description Obtener todos los eventos
// @Tags events
// @Accept json
// @Produce json
// @Param name query string false "Nombre del evento"
// @Param gender query string false "Género del evento"
// @Param location query string false "Ubicación del evento"
// @Success 200 {array} models.Event "Lista de eventos"
// @Failure 400 {object} map[string]string "Formato de parámetros inválido"
// @Failure 401 {object} map[string]string "No autorizado"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /events [get]
// @Security BearerAuth
// GET /events
func (h *EventHandler) GetAllEvents(c *gin.Context) {
	name := c.Query("name")
	gender := c.Query("gender")
	location := c.Query("location")

	events, err := h.service.GetAllEvents(models.EventFilter{Name: name, Gender: gender, Location: location})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch events"})
		return
	}

	c.JSON(http.StatusOK, events)
}

// UpdateEvent godoc
// @Summary Actualizar evento
// @Description Actualizar un evento existente
// @Tags events
// @Accept json
// @Produce json
// @Param id path string true "ID del evento"
// @Param event body models.Event true "Datos del evento"
// @Success 200 {object} map[string]string "Evento actualizado exitosamente"
// @Failure 400 {object} map[string]string "Formato JSON inválido"
// @Failure 401 {object} map[string]string "No autorizado"
// @Failure 404 {object} map[string]string "Evento no encontrado"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /events/{id} [put]
// @Security BearerAuth
// PUT /events/:id
func (h *EventHandler) UpdateEvent(c *gin.Context) {
	id := c.Param("id")
	var event models.Event

	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateEvent(id, &event); err != nil {
		if err.Error() == "Cannot update: event not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event updated successfully"})
}

// DeleteEvent godoc
// @Summary Eliminar evento
// @Description Eliminar un evento existente
// @Tags events
// @Accept json
// @Produce json
// @Param id path string true "ID del evento"
// @Success 200 {object} map[string]string "Evento eliminado exitosamente"
// @Failure 401 {object} map[string]string "No autorizado"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /events/{id} [delete]
// @Security BearerAuth
// DELETE /events/:id
func (h *EventHandler) DeleteEvent(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.DeleteEvent(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete event"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event deleted successfully"})
}

// UpdateAvailabilityForEvent godoc
// @Summary Actualizar disponibilidad para un evento
// @Description Actualizar la disponibilidad para un evento existente
// @Tags events
// @Accept json
// @Produce json
// @Param id path string true "ID del evento"
// @Success 200 {object} map[string]string "Disponibilidad actualizada exitosamente"
// @Failure 401 {object} map[string]string "No autorizado"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /events/availability/{id} [patch]
// @Security BearerAuth
// PATCH /events/availability/:id
func (h *EventHandler) UpdateAvailabilityForEvent(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.UpdateEventAvailability(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update availability for event"})
	    return	
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event availability successfully updated"})
}
