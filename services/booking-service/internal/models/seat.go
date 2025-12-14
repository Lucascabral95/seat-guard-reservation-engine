package models

import (
	"time"

	"gorm.io/gorm"
)

type SeatStatus string

const (
	StatusAvailable SeatStatus = "AVAILABLE"
	StatusLocked    SeatStatus = "LOCKED"
	StatusSold      SeatStatus = "SOLD"
)

type PaymentStatus string

const (
	PaymentPending   PaymentStatus = "PENDING"
	PaymentCompleted PaymentStatus = "COMPLETED"
	PaymentFailed    PaymentStatus = "FAILED"
)

type BaseModel struct {
	ID        string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Event struct {
	BaseModel

	Name        string    `gorm:"not null" json:"name"`
	Description string    `json:"description,omitempty"` // Opcional en JSON
	Location    string    `json:"location"`
	Date        time.Time `json:"date"`
	Price       float64   `gorm:"not null" json:"price"` // Precio base por entrada

	// Relación: Un evento tiene muchos asientos
	Seats []Seat `gorm:"foreignKey:EventID" json:"seats,omitempty"`
}

// --- ASIENTO ---
type Seat struct {
	BaseModel

	// Identificación del asiento (ej: A-10, B-5)
	Section string `json:"section"`                // "VIP", "Platea", "General"
	Number  string `gorm:"not null" json:"number"` // "10", "A1", etc.

	Price float64 `json:"price"` // Precio específico (puede anular el del evento)

	Status SeatStatus `gorm:"type:varchar(20);default:'AVAILABLE'" json:"status"`

	// Control de Bloqueo Temporal
	LockedBy *string    `gorm:"type:text" json:"lockedBy"` // Quien lo bloquea (uuid)
	LockedAt *time.Time `json:"lockedAt"`                  // Cuando se desbloquea

	// Relaciones
	EventID  string  `gorm:"not null" json:"eventId"`
	TicketID *string `json:"ticketId,omitempty"` // ID del ticket final si se vende
}

// --- PAGO / ORDEN ---
type BookingOrder struct {
	BaseModel

	UserID string        `gorm:"not null" json:"userId"`
	Amount float64       `gorm:"not null" json:"amount"`
	Status PaymentStatus `gorm:"default:'PENDING'" json:"status"`

	// Asientos involucrados en esta orden
	//SeatIDs []string `gorm:"type:text[]" json:"seatIds"`
	SeatIDs []string `gorm:"serializer:json" json:"seatIds"`

	// Token o ID de transacción de la pasarela de pago (Stripe/MercadoPago)
	PaymentProviderID string `json:"paymentProviderId,omitempty"`
}
