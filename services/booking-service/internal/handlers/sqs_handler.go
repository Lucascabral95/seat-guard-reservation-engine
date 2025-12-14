package handlers

import (
	"encoding/json"
	"net/http"

	"booking-service/internal/messaging"

	"github.com/gin-gonic/gin"
)

type SQSHandler struct {
	sqs *messaging.SQSClient
}

func NewSQSHandler(sqs *messaging.SQSClient) *SQSHandler {
	return &SQSHandler{sqs: sqs}
}

type StripeWebhookReq struct {
	Object string     `json:"object"`
	Data   StripeData `json:"data"`
}

type StripeData struct {
	Object StripeObject `json:"object"`
}

type StripeObject struct {
	ID            string         `json:"id"`
	PaymentStatus string         `json:"payment_status"`
	Metadata      StripeMetadata `json:"metadata"`
	AmountTotal   float64        `json:"amount_total"`
	PaymentIntent string         `json:"payment_intent"`
}

type StripeMetadata struct {
	UserID  string `json:"user_id"`
	SeatIDs string `json:"seat_ids"`
}

type BookingMessage struct {
	UserID            string  `json:"userId"`
	Amount            float64 `json:"amount"`
	Status            string  `json:"status"`
	SeatIDs           string  `json:"seatIds"`
	PaymentProviderID string  `json:"paymentProviderId"`
}

func (h *SQSHandler) Send(c *gin.Context) {
	var stripeReq StripeWebhookReq
	if err := c.ShouldBindJSON(&stripeReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid webhook json: " + err.Error()})
		return
	}

	internalMsg := BookingMessage{
		UserID:            stripeReq.Data.Object.Metadata.UserID,
		Amount:            stripeReq.Data.Object.AmountTotal / 100,
		SeatIDs:           stripeReq.Data.Object.Metadata.SeatIDs,
		Status:            stripeReq.Data.Object.PaymentStatus,
		PaymentProviderID: stripeReq.Data.Object.PaymentIntent,
	}

	b, err := json.Marshal(internalMsg)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "marshal error: " + err.Error()})
		return
	}

	messageID, err := h.sqs.Send(c.Request.Context(), string(b))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "sqs send error: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"messageId": messageID, "status": "processed"})
}
