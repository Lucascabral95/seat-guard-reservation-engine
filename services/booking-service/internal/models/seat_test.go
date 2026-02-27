package models_test

import (
	"strings"
	"testing"

	"booking-service/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil && strings.Contains(err.Error(), "go-sqlite3 requires cgo") {
		t.Skip("sqlite in-memory requires CGO in this environment")
	}
	require.NoError(t, err)

	err = db.AutoMigrate(
		&models.Event{},
		&models.Seat{},
		&models.BookingOrder{},
		&models.Checkout{},
	)
	require.NoError(t, err)
	return db
}

func TestSeatStatusConstants(t *testing.T) {
	assert.Equal(t, models.SeatStatus("AVAILABLE"), models.StatusAvailable)
	assert.Equal(t, models.SeatStatus("LOCKED"), models.StatusLocked)
	assert.Equal(t, models.SeatStatus("SOLD"), models.StatusSold)
}
