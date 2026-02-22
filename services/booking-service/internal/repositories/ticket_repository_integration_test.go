package repositories

import (
	"booking-service/internal/models"
	"fmt"
	"testing"
	"time"
)

func TestTicketRepository_Integration_CRUD(t *testing.T) {
	db := openIntegrationDB(t)
	repo := NewTicketRepository(db)

	suffix := fmt.Sprintf("%d", time.Now().UnixNano())
	ticketID := "dddddddd-dddd-dddd-dddd-" + suffix[len(suffix)-12:]
	orderID := "eeeeeeee-eeee-eeee-eeee-" + suffix[len(suffix)-12:]

	if err := repo.CreateTicket(nil); err == nil {
		t.Fatalf("expected nil ticket validation error")
	}

	ticket := &models.TicketPDF{
		BaseModel:       models.BaseModel{ID: ticketID},
		PaymentProvider: "STRIPE",
		PaymentIntentID: "pi_repo_test",
		Currency:        "usd",
		Amount:          999,
		Name:            "Repo User",
		Email:           "repo@test.com",
		OrderID:         orderID,
		PDFVersion:      1,
		PDFData:         []byte("pdf"),
	}
	if err := repo.CreateTicket(ticket); err != nil {
		t.Fatalf("create ticket failed: %v", err)
	}

	gotByID, err := repo.FindTicketById(ticketID)
	if err != nil || gotByID.OrderID != orderID {
		t.Fatalf("find by id failed: err=%v ticket=%+v", err, gotByID)
	}

	gotByOrder, err := repo.FindTicketByOrderID(orderID)
	if err != nil || gotByOrder.ID != ticketID {
		t.Fatalf("find by order failed: err=%v ticket=%+v", err, gotByOrder)
	}

	all, err := repo.FindAllTickets()
	if err != nil || len(all) == 0 {
		t.Fatalf("find all failed: err=%v len=%d", err, len(all))
	}

	gotByID.Currency = "ars"
	if err := repo.UpdateTicket(gotByID); err != nil {
		t.Fatalf("update ticket failed: %v", err)
	}

	if err := repo.DeleteTicket(ticketID); err != nil {
		t.Fatalf("delete ticket failed: %v", err)
	}

	if _, err := repo.FindTicketById(ticketID); err == nil {
		t.Fatalf("expected not found after delete")
	}
}
