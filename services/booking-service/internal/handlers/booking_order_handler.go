package handlers

import (
	"booking-service/internal/models"
	"booking-service/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BookingOrderHandler struct {
	service *services.BookingOrderService
}

func NewBookingOrderHandler(service *services.BookingOrderService) *BookingOrderHandler {
	return &BookingOrderHandler{service: service}
}

// CreateBookingOrder godoc
// @Summary Crear booking order
// @Description Crea un nuevo registro de booking order
// @Tags booking-order
// @Accept json
// @Produce json
// @Param bookingOrder body models.BookingOrder true "Datos del booking order"
// @Success 200 {object} models.BookingOrder "Booking order creado exitosamente"
// @Failure 400 {object} map[string]string "Datos inválidos"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /booking-order [post]
// @Security BearerAuth
func (h *BookingOrderHandler) CreateBookingOrder(c *gin.Context) {
	var bookingOrders models.BookingOrder

	if err := c.ShouldBindJSON(&bookingOrders); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateBookingOrder(&bookingOrders); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, bookingOrders)
}

// GetBookingOrders godoc
// @Summary Obtener todos los booking orders
// @Description Obtiene todos los registros de booking orders
// @Tags booking-order
// @Accept json
// @Produce json
// @Success 200 {array} models.BookingOrder "Booking orders encontrados"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /booking-order [get]
// @Security BearerAuth
func (h *BookingOrderHandler) GetBookingOrders(c *gin.Context) {
	bookingOrder, err := h.service.FindAllBookingOrders()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, bookingOrder)
}

type updateBookingOrderReq struct {
	Status            models.PaymentStatus `json:"status" binding:"required"`
	PaymentProviderID string               `json:"paymentProviderId"`
}

// UpdateBookingOrder godoc
// @Summary Actualizar booking order
// @Description Actualiza un registro de booking order
// @Tags booking-order
// @Accept json
// @Produce json
// @Param id path string true "ID del booking order"
// @Param bookingOrder body updateBookingOrderReq true "Datos del booking order"
// @Success 200 {object} models.BookingOrder "Booking order actualizado exitosamente"
// @Failure 400 {object} map[string]string "Datos inválidos"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /booking-order/{id} [put]
// @Security BearerAuth
func (h *BookingOrderHandler) UpdateBookingOrder(c *gin.Context) {
	id := c.Param("id")
	_, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	var req updateBookingOrderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	switch req.Status {
	case models.PaymentPending, models.PaymentCompleted, models.PaymentFailed:
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
		return
	}

	if err := h.service.UpdateBookingOrder(id, req.Status, req.PaymentProviderID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	bookingOrder, err := h.service.FindBookingOrderById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, bookingOrder)
}

// GetBookingOrderById godoc
// @Summary Obtener booking order por ID
// @Description Obtiene un registro de booking order por su ID
// @Tags booking-order
// @Accept json
// @Produce json
// @Param id path string true "ID del booking order"
// @Success 200 {object} models.BookingOrder "Booking order encontrado"
// @Failure 400 {object} map[string]string "ID inválido"
// @Failure 404 {object} map[string]string "Booking order no encontrado"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /booking-order/{id} [get]
// @Security BearerAuth
func (h *BookingOrderHandler) GetBookingOrderById(c *gin.Context) {
	id := c.Param("id")

	_, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	bookingOrder, err := h.service.FindBookingOrderById(id)
	if err != nil {
		if err.Error() == "not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Booking order not found"})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, bookingOrder)
}

// GetAllOrderForUserID godoc
// @Summary Obtener todos los booking orders por ID de usuario
// @Description Obtiene todos los registros de booking orders por ID de usuario
// @Tags booking-order
// @Accept json
// @Produce json
// @Param id path string true "ID del usuario"
// @Success 200 {array} models.BookingOrder "Booking orders encontrados"
// @Failure 400 {object} map[string]string "ID inválido"
// @Failure 404 {object} map[string]string "Booking order no encontrado"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /booking-order/user/{id} [get]
// @Security BearerAuth
func (h *BookingOrderHandler) GetAllOrderForUserID(c *gin.Context) {
	id := c.Param("id")

	_, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	bookingOrders, err := h.service.FindAllOrdersByUserID(id)
	if err != nil {
		if err.Error() == "not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Booking order not found"})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, bookingOrders)
}
