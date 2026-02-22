package seeds

import (
	"booking-service/internal/models"
	"os"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func openSeedIntegrationDB(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := os.Getenv("BOOKING_IT_DATABASE_URL")
	if dsn == "" {
		t.Skip("BOOKING_IT_DATABASE_URL not set; skipping seed integration")
	}
	if os.Getenv("BOOKING_IT_ALLOW_SEED_RESET") != "true" {
		t.Skip("BOOKING_IT_ALLOW_SEED_RESET=true required to run destructive seed test")
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	return db
}

func TestResetAndSeed_Integration(t *testing.T) {
	db := openSeedIntegrationDB(t)

	if err := ResetAndSeed(db); err != nil {
		t.Fatalf("reset and seed failed: %v", err)
	}

	var eventCount int64
	if err := db.Model(&models.Event{}).Count(&eventCount).Error; err != nil {
		t.Fatalf("count events failed: %v", err)
	}
	if eventCount == 0 {
		t.Fatalf("expected seeded events")
	}

	var seatCount int64
	if err := db.Model(&models.Seat{}).Count(&seatCount).Error; err != nil {
		t.Fatalf("count seats failed: %v", err)
	}
	if seatCount == 0 {
		t.Fatalf("expected seeded seats")
	}
}
