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

func (h *CheckoutHandler) GetByOrderID(c *gin.Context) {
	orderID := c.Param("orderID")

	checkout, err := h.service.FindByOrderID(orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, checkout)
}

func (h *CheckoutHandler) GetAll(c *gin.Context) {
	checkouts, err := h.service.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, checkouts)
}

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
