package handlers

import (
	"booking-service/internal/models"
	"booking-service/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CheckoutHandler struct {
	service *services.CheckoutService
}

func NewCheckoutHandler(service *services.CheckoutService) *CheckoutHandler {
	return &CheckoutHandler{service: service}
}

// Create godoc
// @Summary Crear checkout
// @Description Crea un nuevo registro de checkout
// @Tags checkout
// @Accept json
// @Produce json
// @Param checkout body models.Checkout true "Datos del checkout"
// @Success 201 {object} models.Checkout "Checkout creado exitosamente"
// @Failure 400 {object} map[string]string "Datos inválidos"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /checkout [post]
// @Security BearerAuth
func (h *CheckoutHandler) Create(c *gin.Context) {
	var checkout models.Checkout

	if err := c.ShouldBindJSON(&checkout); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if checkout.OrderID == "" ||
		checkout.PaymentIntentID == "" ||
		checkout.Currency == "" ||
		checkout.Amount <= 0 ||
		checkout.CustomerEmail == "" ||
		checkout.CustomerName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Faltan campos obligatorios",
		})
		return
	}

	if checkout.CustomerID != nil && *checkout.CustomerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "customerId inválido",
		})
		return
	}

	if err := h.service.Create(&checkout); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, checkout)
}

// GetByOrderID godoc
// @Summary Obtener checkout por order ID
// @Description Obtiene un registro de checkout por order ID
// @Tags checkout
// @Accept json
// @Produce json
// @Param orderID path string true "Order ID"
// @Success 200 {object} models.Checkout "Checkout encontrado"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /checkout/{orderID} [get]
// @Security BearerAuth
func (h *CheckoutHandler) GetByOrderID(c *gin.Context) {
	orderID := c.Param("orderID")

	checkout, err := h.service.FindByOrderID(orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, checkout)
}

// GetAll godoc
// @Summary Obtener todos los checkouts
// @Description Obtiene todos los registros de checkout
// @Tags checkout
// @Accept json
// @Produce json
// @Success 200 {array} models.Checkout "Checkouts encontrados"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /checkout [get]
// @Security BearerAuth
func (h *CheckoutHandler) GetAll(c *gin.Context) {
	checkouts, err := h.service.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, checkouts)
}

// Update godoc
// @Summary Actualizar checkout
// @Description Actualiza un registro de checkout
// @Tags checkout
// @Accept json
// @Produce json
// @Param checkout body models.Checkout true "Datos del checkout"
// @Success 200 {object} models.Checkout "Checkout actualizado exitosamente"
// @Failure 400 {object} map[string]string "Datos inválidos"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /checkout [put]
// @Security BearerAuth
func (h *CheckoutHandler) Update(c *gin.Context) {
	var checkout models.Checkout

	if err := c.ShouldBindJSON(&checkout); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Update(&checkout); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, checkout)
}
