package repositories

import (
	"booking-service/internal/models"
	"fmt"
	"os"
	"testing"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func openIntegrationDB(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := os.Getenv("BOOKING_IT_DATABASE_URL")
	if dsn == "" {
		t.Skip("BOOKING_IT_DATABASE_URL not set; skipping Postgres integration tests")
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	if err := db.AutoMigrate(&models.Event{}, &models.Seat{}, &models.BookingOrder{}, &models.Checkout{}, &models.TicketPDF{}); err != nil {
		t.Fatalf("failed automigrate: %v", err)
	}
	return db
}

func TestEventAndSeatRepository_Integration_CRUDAndLock(t *testing.T) {
	db := openIntegrationDB(t)
	eventRepo := NewEventRepository(db)
	seatRepo := NewSeatRepository(db)

	suffix := fmt.Sprintf("%d", time.Now().UnixNano())
	eventID := "11111111-1111-1111-1111-" + suffix[len(suffix)-12:]
	seatID := "22222222-2222-2222-2222-" + suffix[len(suffix)-12:]
	userID := "33333333-3333-3333-3333-" + suffix[len(suffix)-12:]

	event := &models.Event{
		BaseModel: models.BaseModel{ID: eventID},
		Name:      "IT Concert " + suffix,
		Location:  "Arena",
		Date:      time.Now().Add(24 * time.Hour),
		Price:     5000,
	}
	if err := eventRepo.Create(event); err != nil {
		t.Fatalf("create event failed: %v", err)
	}

	seat := &models.Seat{
		BaseModel: models.BaseModel{ID: seatID},
		EventID:   eventID,
		Section:   "VIP",
		Number:    "A1",
		Price:     100,
		Status:    models.StatusAvailable,
	}
	if err := seatRepo.Create(seat); err != nil {
		t.Fatalf("create seat failed: %v", err)
	}

	if err := seatRepo.LockSeat(seatID, userID, time.Now().Add(-1*time.Minute)); err != nil {
		t.Fatalf("lock seat failed: %v", err)
	}
	if err := seatRepo.UnlockIfExpired(seatID, time.Now()); err != nil {
		t.Fatalf("unlock seat failed: %v", err)
	}

	gotSeat, err := seatRepo.FindByID(seatID)
	if err != nil {
		t.Fatalf("find seat failed: %v", err)
	}
	if gotSeat.Status != models.StatusAvailable {
		t.Fatalf("expected seat AVAILABLE after unlock, got %s", gotSeat.Status)
	}

	if err := eventRepo.UpdateAvailability(eventID); err != nil {
		t.Fatalf("update availability failed: %v", err)
	}
	gotEvent, err := eventRepo.FindByID(eventID)
	if err != nil {
		t.Fatalf("find event failed: %v", err)
	}
	if gotEvent == nil {
		t.Fatalf("expected event to exist")
	}

	_ = seatRepo.UpdateStatus(seatID, models.StatusSold)
	_ = eventRepo.Delete(eventID)
}
