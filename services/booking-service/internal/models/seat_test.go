package models_test

import (
	"strings"
	"testing"
	"time"

	"booking-service/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type sqliteBaseModel struct {
	ID        string `gorm:"primaryKey;type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type sqliteEvent struct {
	sqliteBaseModel

	Name         string `gorm:"not null"`
	Description  string
	Location     string
	Date         time.Time
	Price        int64 `gorm:"not null"`
	PosterURL    string
	Gender       string `gorm:"type:varchar(20);default:'VARIOS'"`
	Availability string `gorm:"type:varchar(20);default:'HIGH'"`
}

func (sqliteEvent) TableName() string {
	return "events"
}

type sqliteTicketPDF struct {
	sqliteBaseModel

	PaymentProvider string     `gorm:"type:varchar(50);default:'STRIPE'"`
	PaymentIntentID string     `gorm:"type:text;not null"`
	Currency        string     `gorm:"type:varchar(10);not null"`
	Amount          int64      `gorm:"not null"`
	Name            string     `gorm:"type:text;not null"`
	Email           string     `gorm:"type:text;not null"`
	CustomerID      *string    `gorm:"type:varchar(50)"`
	OrderID         string     `gorm:"not null;index"`
	PDFData         []byte     `gorm:"type:blob"`
	PDFGeneratedAt  *time.Time `gorm:"type:timestamp"`
	PDFVersion      int        `gorm:"default:1"`
}

func (sqliteTicketPDF) TableName() string {
	return "ticket_pdfs"
}

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil && strings.Contains(err.Error(), "go-sqlite3 requires cgo") {
		t.Skip("sqlite in-memory requires CGO in this environment")
	}
	require.NoError(t, err)

	err = db.AutoMigrate(&sqliteEvent{})
	require.NoError(t, err)
	return db
}

func TestSeatStatusConstants(t *testing.T) {
	assert.Equal(t, models.SeatStatus("AVAILABLE"), models.StatusAvailable)
	assert.Equal(t, models.SeatStatus("LOCKED"), models.StatusLocked)
	assert.Equal(t, models.SeatStatus("SOLD"), models.StatusSold)
}
