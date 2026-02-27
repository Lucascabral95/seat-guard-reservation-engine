package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"booking-service/internal/messaging"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SQSHandler struct {
	sqs *messaging.SQSClient
}

func NewSQSHandler(sqs *messaging.SQSClient) *SQSHandler {
	return &SQSHandler{sqs: sqs}
}

// Estructuras para parsear el JSON que envía Stripe
type StripeWebhookReq struct {
	ID     string     `json:"id"`
	Type   string     `json:"type"`
	Object string     `json:"object"`
	Data   StripeData `json:"data"`
}

type StripeData struct {
	Object StripeObject `json:"object"`
}

type StripeObject struct {
	ID            string         `json:"id"`             // El PaymentIntent ID o Session ID
	PaymentStatus string         `json:"payment_status"` // "paid", "unpaid"
	Status        string         `json:"status"`         // payment_intent.* usa "succeeded", "requires_payment_method", etc.
	Metadata      StripeMetadata `json:"metadata"`
	AmountTotal   float64        `json:"amount_total"`
	Amount        float64        `json:"amount"`
	PaymentIntent string         `json:"payment_intent"`
}

type StripeMetadata struct {
	UserID  string `json:"user_id"`
	SeatIDs string `json:"seat_ids"`
	EventID string `json:"event_id"`
	OrderID string `json:"order_id"` // <--- ESTO ES LO IMPORTANTE
}

// Mensaje interno que enviamos a SQS y lee la Lambda
type BookingMessage struct {
	UserID            string  `json:"userId"`
	Amount            float64 `json:"amount"`
	Status            string  `json:"status"`
	SeatIDs           string  `json:"seatIds"`
	EventID           string  `json:"eventId"`
	PaymentProviderID string  `json:"paymentProviderId"`
	OrderID           string  `json:"orderId"` // <--- Para que la Lambda sepa qué actualizar
	// Campo para evitar la duplicacion
	Nonce string `json:"nonce"`
	// Observabilidad / idempotencia
	StripeEventID   string `json:"stripeEventId"`
	StripeEventType string `json:"stripeEventType"`
}

// Send Envía un mensaje a SQS con los datos del webhook de Stripe
// @Summary Enviar mensaje a SQS
// @Description Envía un mensaje a SQS con los datos del webhook de Stripe
// @Tags SQS
// @Accept json
// @Produce json
// @Param body body StripeWebhookReq true "Datos del webhook de Stripe"
// @Success 200 {object} map[string]string "Mensaje enviado exitosamente"
// @Failure 400 {object} map[string]string "Solicitud inválida"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /sqs/messaging [post]
func (h *SQSHandler) Send(c *gin.Context) {
	var stripeReq StripeWebhookReq
	if err := c.ShouldBindJSON(&stripeReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid webhook json: " + err.Error()})
		return
	}

	stripeEventID := strings.TrimSpace(stripeReq.ID)
	if stripeEventID == "" {
		stripeEventID = uuid.NewString()
	}

	// Construimos el mensaje interno
	status := strings.TrimSpace(stripeReq.Data.Object.PaymentStatus)
	if status == "" {
		status = strings.TrimSpace(stripeReq.Data.Object.Status)
	}

	amount := stripeReq.Data.Object.AmountTotal
	if amount == 0 {
		amount = stripeReq.Data.Object.Amount
	}

	internalMsg := BookingMessage{
		UserID:            stripeReq.Data.Object.Metadata.UserID,
		Amount:            amount, // Stripe manda centavos, ajusta si guardas float
		SeatIDs:           stripeReq.Data.Object.Metadata.SeatIDs,
		EventID:           stripeReq.Data.Object.Metadata.EventID,
		Status:            status,
		PaymentProviderID: strings.TrimSpace(stripeReq.Data.Object.PaymentIntent),
		OrderID:           stripeReq.Data.Object.Metadata.OrderID, // <--- Pasamos el ID de la orden pendiente
		Nonce:             stripeEventID,
		StripeEventID:     stripeEventID,
		StripeEventType:   strings.TrimSpace(stripeReq.Type),
	}

	if internalMsg.PaymentProviderID == "" {
		internalMsg.PaymentProviderID = strings.TrimSpace(stripeReq.Data.Object.ID)
	}

	if strings.TrimSpace(internalMsg.UserID) == "" || strings.TrimSpace(internalMsg.SeatIDs) == "" || strings.TrimSpace(internalMsg.OrderID) == "" {
		c.JSON(http.StatusOK, gin.H{"status": "ignored", "reason": "missing required metadata"})
		return
	}

	b, err := json.Marshal(internalMsg)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "marshal error: " + err.Error()})
		return
	}

	groupID := strings.TrimSpace(internalMsg.OrderID)
	if groupID == "" {
		groupID = strings.TrimSpace(internalMsg.UserID)
	}

	messageID, err := h.sqs.Send(c.Request.Context(), string(b), groupID, stripeEventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "sqs send error: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"messageId": messageID, "status": "processed"})
}
