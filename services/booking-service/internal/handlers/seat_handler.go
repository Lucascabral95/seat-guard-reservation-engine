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

// CreateSeat Crea un nuevo asiento
// @Summary Crear asiento
// @Description Crear un nuevo asiento
// @Tags Seats
// @Accept json
// @Produce json
// @Param seat body models.Seat true "Datos del asiento"
// @Success 201 {object} models.Seat "Asiento creado satisfactoriamente"
// @Failure 400 {object} map[string]string "Formato JSON inválido"
// @Failure 401 {object} map[string]string "No autorizado"
// @Failure 500 {object} map[string]string "Error al crear el asiento"
// @Router /seats [post]
// @Security BearerAuth
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

// GetSeats Obtiene todos los asientos
// @Summary Obtener asientos
// @Description Obtener todos los asientos
// @Tags Seats
// @Produce json
// @Success 200 {array} models.Seat "Asientos obtenidos satisfactoriamente"
// @Failure 500 {object} map[string]string "Error al obtener los asientos"
// @Router /seats [get]
// @Security BearerAuth
// GET /seats
func (h *SeatHandler) GetSeats(c *gin.Context) {
	seats, err := h.service.GetSeats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch seats"})
		return
	}

	c.JSON(http.StatusOK, seats)
}

// GetSeat Obtiene un asiento por su ID
// @Summary Obtener asiento
// @Description Obtener un asiento por su ID
// @Tags Seats
// @Produce json
// @Param id path string true "ID del asiento"
// @Success 200 {object} models.Seat "Asiento obtenido satisfactoriamente"
// @Failure 400 {object} map[string]string "Formato UUID inválido"
// @Failure 401 {object} map[string]string "No autorizado"
// @Failure 404 {object} map[string]string "Sin resultados"
// @Failure 500 {object} map[string]string "Error al obtener el asiento"
// @Router /seats/{id} [get]
// @Security BearerAuth
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

// GetSeatsByEventId Obtiene todos los asientos por ID de evento
// @Summary Obtener asientos por ID de evento
// @Description Obtener todos los asientos por ID de evento
// @Tags Seats
// @Produce json
// @Param eventId path string true "ID del evento"
// @Success 200 {array} models.Seat "Asientos obtenidos satisfactoriamente"
// @Failure 400 {object} map[string]string "Formato UUID inválido"
// @Failure 401 {object} map[string]string "No autorizado"
// @Failure 404 {object} map[string]string "Asientos no encontrados"
// @Failure 500 {object} map[string]string "Error al obtener los asientos"
// @Router /seats/event/{eventId} [get]
// @Security BearerAuth
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

// UpdateSeat Actualiza el estado de un asiento
// @Summary Actualizar estado de asiento
// @Description Actualizar el estado de un asiento
// @Tags Seats
// @Accept json
// @Produce json
// @Param id path string true "ID del asiento"
// @Param status body string true "Nuevo estado del asiento"
// @Success 200 {object} map[string]string "Asiento actualizado satisfactoriamente"
// @Failure 400 {object} map[string]string "Formato UUID inválido"
// @Failure 401 {object} map[string]string "No autorizado"
// @Failure 404 {object} map[string]string "Asiento no encontrado"
// @Failure 500 {object} map[string]string "Error al actualizar el asiento"
// @Router /seats/{id} [patch]
// @Security BearerAuth
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

// LockSeat Bloquear por 15 minutos. Requisitos: recibir uid y id del asiento
// @Summary Bloquear asiento
// @Description Bloquear asiento por 15 minutos
// @Tags Seats
// @Accept json
// @Produce json
// @Param id path string true "ID del asiento"
// @Param uid path string true "UID del usuario"
// @Success 200 {object} map[string]string "Asiento bloqueado satisfactoriamente"
// @Failure 400 {object} map[string]string "Formato UUID inválido"
// @Failure 401 {object} map[string]string "No autorizado"
// @Failure 404 {object} map[string]string "Asiento no encontrado"
// @Failure 500 {object} map[string]string "Error al bloquear el asiento"
// @Router /seats/lock/{id}/uid/{uid} [patch]
// @Security BearerAuth
// PATCH /seats/lock/:id/uid/:uid
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
