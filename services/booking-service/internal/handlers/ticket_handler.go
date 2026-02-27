package handlers

import (
	"booking-service/internal/models"
	"booking-service/internal/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TicketHandler struct {
	ticketService       *services.TicketService
	pdfService          *services.PDFService
	bookingOrderService *services.BookingOrderService
	checkoutService     *services.CheckoutService
}

func NewTicketHandler(
	ticketService *services.TicketService,
	pdfService *services.PDFService,
	bookingOrderService *services.BookingOrderService,
	checkoutService *services.CheckoutService,
) *TicketHandler {
	return &TicketHandler{
		ticketService:       ticketService,
		pdfService:          pdfService,
		bookingOrderService: bookingOrderService,
		checkoutService:     checkoutService,
	}
}

// GetTicketMetadata godoc
// @Summary Obtener metadata del ticket
// @Description Obtiene los datos del ticket sin el PDF binario
// @Tags tickets
// @Accept json
// @Produce json
// @Param orderID path string true "ID del order"
// @Success 200 {object} map[string]interface{} "Metadata del ticket"
// @Failure 400 {object} map[string]string "ID inválido"
// @Failure 401 {object} map[string]string "No autenticado"
// @Failure 403 {object} map[string]string "Acceso denegado"
// @Failure 404 {object} map[string]string "Ticket no encontrado"
// @Router /tickets/{orderID} [get]
// @Security BearerAuth
// GetTicketMetadata obtiene los datos del ticket sin el PDF binario
// GET /api/v1/tickets/:orderID
func (h *TicketHandler) GetTicketMetadata(c *gin.Context) {
	orderID := c.Param("orderID")

	userIDRaw, exists := c.Get("userID")
	if !exists || userIDRaw == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userID, ok := userIDRaw.(string)
	if !ok || userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user identity"})
		return
	}

	ticket, err := h.ticketService.GetTicketByOrderID(orderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
		return
	}

	if err := h.ticketService.ValidateTicketOwnership(ticket.ID, userID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":              ticket.ID,
		"orderID":         ticket.OrderID,
		"eventName":       ticket.EventName,
		"eventHour":       ticket.EventHour,
		"customerName":    ticket.Name,
		"customerEmail":   ticket.Email,
		"amount":          ticket.Amount,
		"currency":        ticket.Currency,
		"paymentProvider": ticket.PaymentProvider,
		"pdfGenerated":    ticket.PDFGeneratedAt != nil,
		"pdfVersion":      ticket.PDFVersion,
		"seats":           ticket.Items,
		"createdAt":       ticket.CreatedAt,
	})
}

// DownloadTicketPDF godoc
// @Summary Descargar PDF del ticket
// @Description Genera y descarga el PDF del ticket
// @Tags tickets
// @Accept json
// @Produce application/pdf
// @Param orderID path string true "ID del order"
// @Success 200 {file} file "PDF del ticket"
// @Failure 400 {object} map[string]string "ID inválido"
// @Failure 401 {object} map[string]string "No autenticado"
// @Failure 403 {object} map[string]string "Acceso denegado"
// @Failure 404 {object} map[string]string "Ticket no encontrado"
// @Router /tickets/{orderID}/download [get]
// @Security BearerAuth
// DownloadTicketPDF genera y descarga el PDF del ticket
// GET /api/v1/tickets/:orderID/download
func (h *TicketHandler) DownloadTicketPDF(c *gin.Context) {
	orderID := c.Param("orderID")

	// Esto es para pedir Bearer token (JWT)
	ticket, err := h.ticketService.GetTicketByOrderID(orderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
		return
	}

	if ticket.PDFData != nil && len(ticket.PDFData) > 0 {
		filename := fmt.Sprintf("ticket-%s.pdf", ticket.OrderID[:8])
		c.Header("Content-Type", "application/pdf")
		c.Header("Content-Disposition", fmt.Sprintf("inline; filename=%s", filename))
		c.Data(http.StatusOK, "application/pdf", ticket.PDFData)
		return
	}

	pdfBytes, err := h.pdfService.GenerateTicket(ticket)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate PDF: " + err.Error(),
		})
		return
	}

	if err := h.ticketService.UpdateTicketPDF(ticket.ID, pdfBytes); err != nil {
		fmt.Printf("⚠️ Warning: Failed to save PDF to database: %v\n", err)
	}

	filename := fmt.Sprintf("ticket-%s.pdf", ticket.OrderID[:8])
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", fmt.Sprintf("inline; filename=%s", filename))
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}

// RegenerateTicketPDF godoc
// @Summary Regenerar PDF del ticket
// @Description Regenera el PDF del ticket
// @Tags tickets
// @Accept json
// @Produce application/pdf
// @Param orderID path string true "ID del order"
// @Success 200 {file} file "PDF del ticket"
// @Failure 400 {object} map[string]string "ID inválido"
// @Failure 401 {object} map[string]string "No autenticado"
// @Failure 403 {object} map[string]string "Acceso denegado"
// @Failure 404 {object} map[string]string "Ticket no encontrado"
// @Router /tickets/{orderID}/regenerate [post]
// @Security BearerAuth
// RegenerateTicketPDF regenera el PDF de un ticket
// POST /api/v1/tickets/:orderID/regenerate
func (h *TicketHandler) RegenerateTicketPDF(c *gin.Context) {
	orderID := c.Param("orderID")

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// 1. Obtener ticket
	ticket, err := h.ticketService.GetTicketByOrderID(orderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
		return
	}

	// 2. Validar ownership
	if err := h.ticketService.ValidateTicketOwnership(ticket.ID, userID.(string)); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// 3. Generar nuevo PDF
	pdfBytes, err := h.pdfService.GenerateTicket(ticket)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate PDF: " + err.Error(),
		})
		return
	}

	// 4. Actualizar en DB
	if err := h.ticketService.UpdateTicketPDF(ticket.ID, pdfBytes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save PDF: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "PDF regenerated successfully",
		"ticketID":   ticket.ID,
		"pdfVersion": ticket.PDFVersion + 1,
	})
}

// GetAllTickets godoc
// @Summary Obtener todos los tickets
// @Description Obtiene todos los tickets (solo admin)
// @Tags tickets
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Lista de tickets"
// @Failure 500 {object} map[string]string "Error interno"
// @Router /tickets [get]
// @Security BearerAuth
// GetAllTickets obtiene todos los tickets (admin)
// GET /api/v1/tickets
func (h *TicketHandler) GetAllTickets(c *gin.Context) {
	tickets, err := h.ticketService.GetAllTickets()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch tickets: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total":   len(tickets),
		"tickets": tickets,
	})
}

// GetTicketByID godoc
// @Summary Obtener ticket por ID
// @Description Obtiene un ticket específico por su ID (solo owner o admin)
// @Tags tickets
// @Accept json
// @Produce json
// @Param ticketID path string true "ID del ticket"
// @Success 200 {object} models.TicketPDF "Ticket encontrado"
// @Failure 404 {object} map[string]string "Ticket no encontrado"
// @Failure 403 {object} map[string]string "Acceso denegado"
// @Failure 500 {object} map[string]string "Error interno"
// @Router /tickets/by-id/{ticketID} [get]
// @Security BearerAuth
// GetTicketByID obtiene un ticket por su ID
// GET /api/v1/tickets/by-id/:ticketID
func (h *TicketHandler) GetTicketByID(c *gin.Context) {
	ticketID := c.Param("ticketID")

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	ticket, err := h.ticketService.GetTicketByID(ticketID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
		return
	}

	// Validar ownership
	if err := h.ticketService.ValidateTicketOwnership(ticket.ID, userID.(string)); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	c.JSON(http.StatusOK, ticket)
}

// DeleteTicket godoc
// @Summary Eliminar ticket
// @Description Elimina un ticket (soft delete) - solo owner o admin
// @Tags tickets
// @Accept json
// @Produce json
// @Param orderID path string true "ID del order"
// @Success 200 {object} map[string]string "Ticket eliminado"
// @Failure 404 {object} map[string]string "Ticket no encontrado"
// @Failure 403 {object} map[string]string "Acceso denegado"
// @Failure 500 {object} map[string]string "Error interno"
// @Router /tickets/{orderID} [delete]
// @Security BearerAuth
// DeleteTicket elimina un ticket (soft delete)
// DELETE /api/v1/tickets/:orderID
func (h *TicketHandler) DeleteTicket(c *gin.Context) {
	orderID := c.Param("orderID")

	// Obtener userID del JWT middleware (solo admin debería poder hacer esto)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Obtener ticket
	ticket, err := h.ticketService.GetTicketByOrderID(orderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
		return
	}

	// Validar ownership
	if err := h.ticketService.ValidateTicketOwnership(ticket.ID, userID.(string)); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Eliminar
	if err := h.ticketService.DeleteTicket(ticket.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete ticket: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Ticket deleted successfully",
		"ticketID": ticket.ID,
	})
}

// CreateTicketFromEndpoint godoc
// @Summary Crear ticket desde endpoint
// @Description Crea un ticket a partir de un order completado
// @Tags tickets
// @Accept json
// @Produce json
// @Param request body object true "Datos del ticket"
// @Success 201 {object} models.TicketPDF "Ticket creado"
// @Failure 400 {object} map[string]string "Datos inválidos"
// @Failure 404 {object} map[string]string "Order no encontrado"
// @Failure 409 {object} map[string]string "Ticket ya existe"
// @Failure 500 {object} map[string]string "Error interno"
// @Router /tickets [post]
// @Security BearerAuth
// Creacion de ticket desde endpoint
// ✅ CÓDIGO CORRECTO
func (h *TicketHandler) CreateTicketFromEndpoint(c *gin.Context) {
	// 1. Recibir solo el orderID
	var req struct {
		OrderID string `json:"orderId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. Buscar orden en DB
	order, err := h.bookingOrderService.FindBookingOrderById(req.OrderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	// 3. Verificar que esté COMPLETED
	if order.Status != models.PaymentCompleted {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order not completed yet"})
		return
	}

	// 4. Buscar checkout en DB
	checkout, err := h.checkoutService.FindByOrderID(req.OrderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Checkout not found"})
		return
	}

	// 5. Crear ticket con datos reales de DB
	ticket, err := h.ticketService.CreateTicketFromOrder(checkout, order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Ticket created successfully",
		"ticketId": ticket.ID,
		"orderId":  ticket.OrderID,
	})
}
