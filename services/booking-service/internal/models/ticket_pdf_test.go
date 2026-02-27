package models_test

import (
	"testing"

	"booking-service/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestTicketPDF_TableName(t *testing.T) {
	ticket := models.TicketPDF{}
	assert.Equal(t, "ticket_pdfs", ticket.TableName())
}

func TestTicketPDF_CustomerIDOptional(t *testing.T) {
	withoutCustomer := models.TicketPDF{}
	assert.Nil(t, withoutCustomer.CustomerID)

	id := "cus_123"
	withCustomer := models.TicketPDF{CustomerID: &id}
	assert.NotNil(t, withCustomer.CustomerID)
	assert.Equal(t, "cus_123", *withCustomer.CustomerID)
}

func TestTicketPDF_IgnoredFieldsMapping(t *testing.T) {
	ticket := models.TicketPDF{
		EventName: "Rock in Rio",
		EventHour: "21:00",
		Items: []models.Seat{
			{Number: "A1"},
		},
	}

	assert.Equal(t, "Rock in Rio", ticket.EventName)
	assert.Equal(t, "21:00", ticket.EventHour)
	assert.Len(t, ticket.Items, 1)
	assert.Equal(t, "A1", ticket.Items[0].Number)
}
