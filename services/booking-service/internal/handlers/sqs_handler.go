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

type SendMessageReq struct {
	UserId            string   `json:"userId" binding:"required"`
	Amount            float64  `json:"amount" binding:"required"`
	SeatIds           []string `json:"seatIds" binding:"required"`
	PaymentProviderId string   `json:"paymentProviderId" binding:"required"`
}

func (h *SQSHandler) Send(c *gin.Context) {
	var req SendMessageReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	b, err := json.Marshal(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	messageID, err := h.sqs.Send(c.Request.Context(), string(b))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"messageId": messageID})
}
