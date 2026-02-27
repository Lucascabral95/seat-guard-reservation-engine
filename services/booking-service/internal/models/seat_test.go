package models_test

import (
	"testing"

	"booking-service/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestSeatStatusConstants(t *testing.T) {
	assert.Equal(t, models.SeatStatus("AVAILABLE"), models.StatusAvailable)
	assert.Equal(t, models.SeatStatus("LOCKED"), models.StatusLocked)
	assert.Equal(t, models.SeatStatus("SOLD"), models.StatusSold)
}

func TestPaymentStatusConstants(t *testing.T) {
	assert.Equal(t, models.PaymentStatus("PENDING"), models.PaymentPending)
	assert.Equal(t, models.PaymentStatus("COMPLETED"), models.PaymentCompleted)
	assert.Equal(t, models.PaymentStatus("FAILED"), models.PaymentFailed)
}

func TestSeatStructMapping(t *testing.T) {
	seat := models.Seat{
		Section: "VIP",
		Number:  "A1",
		Price:   100.5,
		Status:  models.StatusAvailable,
		EventID: "event-123",
	}

	assert.Equal(t, "VIP", seat.Section)
	assert.Equal(t, "A1", seat.Number)
	assert.Equal(t, 100.5, seat.Price)
	assert.Equal(t, models.StatusAvailable, seat.Status)
	assert.Equal(t, "event-123", seat.EventID)
}
