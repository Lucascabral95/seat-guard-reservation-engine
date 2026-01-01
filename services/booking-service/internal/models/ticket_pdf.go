package models

import "time"

type TicketPDF struct {
	BaseModel

	PaymentProvider string `gorm:"type:varchar(50);default:'STRIPE'" json:"paymentProvider"`
	PaymentIntentID string `gorm:"type:text;not null" json:"paymentIntentId"`
	Currency        string `gorm:"type:varchar(10);not null" json:"currency"`
	Amount          int64  `gorm:"not null" json:"amount"`

	Name       string  `gorm:"type:text;not null" json:"name"`
	Email      string  `gorm:"type:text;not null" json:"email"`
	CustomerID *string `gorm:"type:varchar(50)" json:"customerId,omitempty"`

	OrderID string `gorm:"not null;index" json:"orderId"`

	EventName string `gorm:"-" json:"eventName,omitempty"`
	EventHour string `gorm:"-" json:"eventHour,omitempty"`

	Items []Seat `gorm:"-" json:"items,omitempty"`

	// ✅ NUEVO: PDF binario y metadata
	PDFData        []byte     `gorm:"type:bytea" json:"-"`
	PDFGeneratedAt *time.Time `gorm:"type:timestamp" json:"pdfGeneratedAt,omitempty"`
	PDFVersion     int        `gorm:"default:1" json:"pdfVersion"`
}

func (TicketPDF) TableName() string {
	return "ticket_pdfs"
}
