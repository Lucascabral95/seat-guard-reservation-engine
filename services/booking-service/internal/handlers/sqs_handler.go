package handlers

import (
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
	Message string `json:"message" binding:"required"`
}

func (h *SQSHandler) Send(c *gin.Context) {
	var req SendMessageReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	messageID, err := h.sqs.Send(c.Request.Context(), req.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"messageId": messageID})
}
