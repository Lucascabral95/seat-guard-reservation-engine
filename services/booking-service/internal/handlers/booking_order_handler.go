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
