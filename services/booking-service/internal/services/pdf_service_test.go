package services

import (
	"booking-service/internal/models"
	"bytes"
	"testing"
)

func TestPDFService_GenerateTicket(t *testing.T) {
	svc := NewPDFService()

	if _, err := svc.GenerateTicket(nil); err == nil {
		t.Fatalf("expected error for nil ticket")
	}

	pdf, err := svc.GenerateTicket(&models.TicketPDF{
		PaymentProvider: "STRIPE",
		PaymentIntentID: "pi_1",
		Currency:        "USD",
		Amount:          12345,
		Name:            "Test User",
		Email:           "test@example.com",
		OrderID:         "12345678-1234-1234-1234-123456789012",
		EventName:       "Rock Fest",
		EventHour:       "21:00",
		PDFVersion:      1,
	})
	if err != nil {
		t.Fatalf("unexpected error generating PDF: %v", err)
	}
	if len(pdf) == 0 {
		t.Fatalf("expected non-empty PDF bytes")
	}
	if !bytes.HasPrefix(pdf, []byte("%PDF")) {
		t.Fatalf("expected PDF header, got prefix %q", string(pdf[:4]))
	}

}
