package models_test

import (
	"testing"
	"time"

	"booking-service/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupTestDBWithTicket(t *testing.T) *gorm.DB {
    db := setupTestDB(t)
    require.NoError(t, db.AutoMigrate(&models.TicketPDF{}))
    return db
}

func TestTicketPDF_Create(t *testing.T) {
    db := setupTestDBWithTicket(t)

    ticket := models.TicketPDF{
        BaseModel:       models.BaseModel{ID: uuid.NewString()},
        PaymentIntentID: "pi_test_123456",
        Currency:        "ARS",
        Amount:          75000,
        Name:            "Juan Pérez",
        Email:           "juan@example.com",
        OrderID:         uuid.NewString(),
    }

    result := db.Create(&ticket)
    assert.NoError(t, result.Error)
    assert.Equal(t, int64(1), result.RowsAffected)
}

func TestTicketPDF_RequiredFields(t *testing.T) {
    db := setupTestDBWithTicket(t)

    tests := []struct {
        name   string
        ticket models.TicketPDF
    }{
        {
            name: "missing PaymentIntentID",
            ticket: models.TicketPDF{
                BaseModel: models.BaseModel{ID: uuid.NewString()},
                Currency:  "ARS",
                Amount:    50000,
                Name:      "Ana",
                Email:     "ana@test.com",
                OrderID:   uuid.NewString(),
            },
        },
        {
            name: "missing Name",
            ticket: models.TicketPDF{
                BaseModel:       models.BaseModel{ID: uuid.NewString()},
                PaymentIntentID: "pi_test_abc",
                Currency:        "ARS",
                Amount:          50000,
                Email:   "ana@test.com",
                OrderID: uuid.NewString(),
            },
        },
        {
            name: "missing Email",
            ticket: models.TicketPDF{
                BaseModel:       models.BaseModel{ID: uuid.NewString()},
                PaymentIntentID: "pi_test_abc",
                Currency:        "ARS",
                Amount:          50000,
                Name:            "Ana",
                OrderID: uuid.NewString(),
            },
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            result := db.Create(&tc.ticket)
            assert.Error(t, result.Error, "esperaba error en: %s", tc.name)
        })
    }
}

func TestTicketPDF_DefaultValues(t *testing.T) {
    db := setupTestDBWithTicket(t)

    ticket := models.TicketPDF{
        BaseModel:       models.BaseModel{ID: uuid.NewString()},
        PaymentIntentID: "pi_test_defaults",
        Currency:        "USD",
        Amount:          100,
        Name:            "Carlos Test",
        Email:           "carlos@test.com",
        OrderID:         uuid.NewString(),
    }

    db.Create(&ticket)

    var found models.TicketPDF
    db.First(&found, "id = ?", ticket.ID)

    assert.Equal(t, "STRIPE", found.PaymentProvider)
    assert.Equal(t, 1, found.PDFVersion)
    assert.Nil(t, found.PDFGeneratedAt)
}

func TestTicketPDF_PDFData_StoreAndRetrieve(t *testing.T) {
    db := setupTestDBWithTicket(t)

    now := time.Now()
    fakePDF := []byte("%PDF-1.4 fake content for testing")

    ticket := models.TicketPDF{
        BaseModel:       models.BaseModel{ID: uuid.NewString()},
        PaymentIntentID: "pi_test_pdf",
        Currency:        "ARS",
        Amount:          60000,
        Name:            "María García",
        Email:           "maria@test.com",
        OrderID:         uuid.NewString(),
        PDFData:         fakePDF,
        PDFGeneratedAt:  &now,
        PDFVersion:      2,
    }

    db.Create(&ticket)

    var found models.TicketPDF
    db.First(&found, "id = ?", ticket.ID)

    assert.Equal(t, fakePDF, found.PDFData)
    assert.NotNil(t, found.PDFGeneratedAt)
    assert.Equal(t, 2, found.PDFVersion)
}

func TestTicketPDF_OptionalCustomerID(t *testing.T) {
    db := setupTestDBWithTicket(t)

    ticketWithout := models.TicketPDF{
        BaseModel:       models.BaseModel{ID: uuid.NewString()},
        PaymentIntentID: "pi_no_customer",
        Currency:        "ARS",
        Amount:          30000,
        Name:            "Guest User",
        Email:           "guest@test.com",
        OrderID:         uuid.NewString(),
    }
    require.NoError(t, db.Create(&ticketWithout).Error)

    customerID := "cus_stripe_abc123"
    ticketWith := models.TicketPDF{
        BaseModel:       models.BaseModel{ID: uuid.NewString()},
        PaymentIntentID: "pi_with_customer",
        Currency:        "ARS",
        Amount:          30000,
        Name:            "Registered User",
        Email:           "reg@test.com",
        OrderID:         uuid.NewString(),
        CustomerID:      &customerID,
    }
    require.NoError(t, db.Create(&ticketWith).Error)

    var found models.TicketPDF
    db.First(&found, "id = ?", ticketWithout.ID)
    assert.Nil(t, found.CustomerID)

    db.First(&found, "id = ?", ticketWith.ID)
    assert.Equal(t, customerID, *found.CustomerID)
}

func TestTicketPDF_GormIgnoredFields(t *testing.T) {
    db := setupTestDBWithTicket(t)

    ticket := models.TicketPDF{
        BaseModel:       models.BaseModel{ID: uuid.NewString()},
        PaymentIntentID: "pi_ignored",
        Currency:        "ARS",
        Amount:          10000,
        Name:            "Test Ignored",
        Email:           "ignored@test.com",
        OrderID:         uuid.NewString(),
        EventName:       "Este valor no se guarda",
        EventHour:       "21:00",
    }

    db.Create(&ticket)

    var found models.TicketPDF
    db.First(&found, "id = ?", ticket.ID)

    assert.Empty(t, found.EventName)
    assert.Empty(t, found.EventHour)
    assert.Nil(t, found.Items)
}

func TestTicketPDF_TableName(t *testing.T) {
    ticket := models.TicketPDF{}
    assert.Equal(t, "ticket_pdfs", ticket.TableName())
}
