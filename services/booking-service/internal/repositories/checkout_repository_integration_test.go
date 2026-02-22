package repositories

import (
	"booking-service/internal/models"
	"fmt"
	"testing"
	"time"
)

func TestCheckoutRepository_Integration_FindByOrderIDAndUpdate(t *testing.T) {
	db := openIntegrationDB(t)
	checkoutRepo := NewCheckoutRepository(db)
	orderRepo := NewBookingOrderRepository(db)
	eventRepo := NewEventRepository(db)
	seatRepo := NewSeatRepository(db)

	suffix := fmt.Sprintf("%d", time.Now().UnixNano())
	eventID := "99999999-9999-9999-9999-" + suffix[len(suffix)-12:]
	seatID := "aaaaaaaa-aaaa-aaaa-aaaa-" + suffix[len(suffix)-12:]
	orderID := "bbbbbbbb-bbbb-bbbb-bbbb-" + suffix[len(suffix)-12:]
	checkoutID := "cccccccc-cccc-cccc-cccc-" + suffix[len(suffix)-12:]

	_ = eventRepo.Create(&models.Event{BaseModel: models.BaseModel{ID: eventID}, Name: "Checkout Event", Location: "Arena", Date: time.Now().Add(24 * time.Hour), Price: 1200})
	_ = seatRepo.Create(&models.Seat{BaseModel: models.BaseModel{ID: seatID}, EventID: eventID, Section: "B", Number: "10", Price: 12, Status: models.StatusAvailable})
	_ = orderRepo.Create(&models.BookingOrder{BaseModel: models.BaseModel{ID: orderID}, UserID: "u1", Amount: 1200, Status: models.PaymentPending, SeatIDs: []string{seatID}})

	checkout := &models.Checkout{
		BaseModel:       models.BaseModel{ID: checkoutID},
		OrderID:         orderID,
		PaymentProvider: "STRIPE",
		PaymentIntentID: "pi_test",
		Currency:        "usd",
		Amount:          1200,
		CustomerEmail:   "buyer@test.com",
		CustomerName:    "Buyer",
	}
	if err := checkoutRepo.Create(checkout); err != nil {
		t.Fatalf("create checkout failed: %v", err)
	}

	got, err := checkoutRepo.FindByOrderID(orderID)
	if err != nil {
		t.Fatalf("find by order id failed: %v", err)
	}
	if got.Order.ID != orderID || len(got.Order.Items) != 1 {
		t.Fatalf("expected preloaded order with seats, got %+v", got.Order)
	}

	got.Currency = "ars"
	if err := checkoutRepo.Update(got); err != nil {
		t.Fatalf("update checkout failed: %v", err)
	}

	all, err := checkoutRepo.FindAll()
	if err != nil || len(all) == 0 {
		t.Fatalf("expected checkout list, err=%v len=%d", err, len(all))
	}
}
