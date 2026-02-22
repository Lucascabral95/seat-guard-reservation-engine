package database

import (
	"booking-service/internal/config"
	"os"
	"testing"
	"time"
)

func TestDatabase_Integration_InitMigrateClose(t *testing.T) {
	dsn := os.Getenv("BOOKING_IT_DATABASE_URL")
	if dsn == "" {
		t.Skip("BOOKING_IT_DATABASE_URL not set; skipping DB integration")
	}

	cfg := &config.Config{
		DBUrl:             dsn,
		DbMaxOpenConns:    5,
		DbMaxIdleConns:    2,
		DbConnMaxLifeTime: 2 * time.Minute,
		DbConnMaxIdleTime: 1 * time.Minute,
	}

	db := InitDB(t.Context(), cfg)
	if db == nil {
		t.Fatalf("expected non-nil db")
	}

	if err := RunMigrations(db); err != nil {
		t.Fatalf("run migrations failed: %v", err)
	}

	if err := CloseDB(db); err != nil {
		t.Fatalf("close db failed: %v", err)
	}
}
