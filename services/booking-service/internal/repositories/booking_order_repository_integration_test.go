package repositories

import (
	"booking-service/internal/models"
	"fmt"
	"testing"
	"time"
)

func TestBookingOrderRepository_Integration_FindAndUpdate(t *testing.T) {
	db := openIntegrationDB(t)
	repo := NewBookingOrderRepository(db)
	eventRepo := NewEventRepository(db)
	seatRepo := NewSeatRepository(db)

	suffix := fmt.Sprintf("%d", time.Now().UnixNano())
	eventID := "44444444-4444-4444-4444-" + suffix[len(suffix)-12:]
	seat1ID := "55555555-5555-5555-5555-" + suffix[len(suffix)-12:]
	seat2ID := "66666666-6666-6666-6666-" + suffix[len(suffix)-12:]
	orderID := "77777777-7777-7777-7777-" + suffix[len(suffix)-12:]
	userID := "88888888-8888-8888-8888-" + suffix[len(suffix)-12:]

	_ = eventRepo.Create(&models.Event{BaseModel: models.BaseModel{ID: eventID}, Name: "Repo Test Event", Location: "Arena", Date: time.Now().Add(24 * time.Hour), Price: 1000})
	_ = seatRepo.Create(&models.Seat{BaseModel: models.BaseModel{ID: seat1ID}, EventID: eventID, Section: "A", Number: "1", Price: 10, Status: models.StatusAvailable})
	_ = seatRepo.Create(&models.Seat{BaseModel: models.BaseModel{ID: seat2ID}, EventID: eventID, Section: "A", Number: "2", Price: 20, Status: models.StatusAvailable})

	order := &models.BookingOrder{
		BaseModel: models.BaseModel{ID: orderID},
		UserID:    userID,
		Amount:    3000,
		Status:    models.PaymentPending,
		SeatIDs:   []string{seat2ID, seat1ID},
	}
	if err := repo.Create(order); err != nil {
		t.Fatalf("create order failed: %v", err)
	}

	got, err := repo.FindByID(orderID)
	if err != nil {
		t.Fatalf("find by id failed: %v", err)
	}
	if len(got.Items) != 2 || got.Items[0].ID != seat2ID || got.Items[1].ID != seat1ID {
		t.Fatalf("expected ordered items from seat IDs, got %+v", got.Items)
	}

	if err := repo.UpdateStatus(orderID, models.PaymentFailed); err != nil {
		t.Fatalf("update status failed: %v", err)
	}
	if err := repo.Update(orderID, models.PaymentCompleted, "pi_123"); err != nil {
		t.Fatalf("update payment provider failed: %v", err)
	}

	byUser, err := repo.FindAllOrdersByUserID(userID)
	if err != nil || len(byUser) == 0 {
		t.Fatalf("expected orders by user, err=%v len=%d", err, len(byUser))
	}
}
