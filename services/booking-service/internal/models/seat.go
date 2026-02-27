package models

import (
	"time"

	"gorm.io/gorm"
)

// @Description Estado del asiento
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

type Availability string

const (
	AvailabilityHigh    Availability = "HIGH"
	AvailabilityMedium  Availability = "MEDIUM"
	AvailabilityLow     Availability = "LOW"
	AvailabilitySoldOut Availability = "SOLD_OUT"
)

type Gender string

const (
	Electronica Gender = "ELECTRONICA"
	Rock        Gender = "ROCK"
	Pop         Gender = "POP"
	Jazz        Gender = "JAZZ"
	Teatro      Gender = "TEATRO"
	Varios      Gender = "VARIOS"
	Metal       Gender = "METAL"
)

// BaseModel es la estructura base para todas las entidades
type BaseModel struct {
	ID        string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Event representa un evento musical o de entretenimiento
type Event struct {
	BaseModel

	Name         string       `gorm:"not null" json:"name"`
	Description  string       `json:"description,omitempty"` // Opcional en JSON
	Location     string       `json:"location"`
	Date         time.Time    `json:"date"`
	Price        int64        `gorm:"not null" json:"price"` // Precio base por entrada
	PosterURL    string       `json:"posterUrl"`
	Gender       string       `gorm:"type:varchar(20);default:'VARIOS'" json:"gender"`
	// Disponibilidad del evento
	// enums: HIGH, MEDIUM, LOW, SOLD_OUT
	Availability Availability `gorm:"type:varchar(20);default:'HIGH'" json:"availability"`

	Seats []Seat `gorm:"foreignKey:EventID" json:"seats,omitempty"`
}

// Seat representa un asiento específico en un evento
type Seat struct {
	BaseModel

	// Identificación del asiento (ej: A-10, B-5)
	Section string `json:"section"`                // "VIP", "Platea", "General"
	Number  string `gorm:"not null" json:"number"` // "10", "A1", etc.

	Price float64 `json:"price"` 

	// Estado actual del asiento
	// enums: AVAILABLE, RESERVED, SOLD, BLOCKED
	Status SeatStatus `gorm:"type:varchar(20);default:'AVAILABLE'" json:"status"`

	// Control de Bloqueo Temporal
	LockedBy *string    `gorm:"type:text" json:"lockedBy"` // Quien lo bloquea (uuid)
	LockedAt *time.Time `json:"lockedAt"`                  // Cuando se desbloquea

	// Relaciones
	EventID  string  `gorm:"not null" json:"eventId"`
	TicketID *string `json:"ticketId,omitempty"` // ID del ticket final si se vende

	EventName string `gorm:"-" json:"eventName,omitempty"`
	EventHour string `gorm:"-" json:"eventHour,omitempty"`
}

// BookingOrder representa una orden de pago para uno o más asientos
type BookingOrder struct {
	BaseModel

	UserID string        `gorm:"not null" json:"userId"`
	Amount int64         `gorm:"not null" json:"amount"`
	Status PaymentStatus `gorm:"default:'PENDING'" json:"status"`

	// Asientos involucrados en esta orden
	//SeatIDs []string `gorm:"type:text[]" json:"seatIds"`
	SeatIDs []string `gorm:"serializer:json" json:"seatIds"`
	Items   []Seat   `gorm:"-" json:"items,omitempty"`

	// Token o ID de transacción de la pasarela de pago (Stripe/MercadoPago)
	PaymentProviderID string `json:"paymentProviderId,omitempty"`
	EventName         string `gorm:"-" json:"eventName,omitempty"`
	EventHour         string `gorm:"-" json:"eventHour,omitempty"`
}

// Checkout representa los datos del usuario pagador
// Datos del usuario pagador
type Checkout struct {
	BaseModel

	OrderID string       `gorm:"not null;index" json:"orderId"`
	Order   BookingOrder `gorm:"foreignKey:OrderID" json:"order,omitempty"`

	PaymentProvider string `gorm:"type:varchar(50);default:'STRIPE'" json:"paymentProvider"`
	PaymentIntentID string `gorm:"type:text;not null" json:"paymentIntentId"`

	Currency string `gorm:"type:varchar(10);not null" json:"currency"`
	Amount   int64  `gorm:"not null" json:"amount"`

	CustomerEmail string `gorm:"type:text;not null" json:"customerEmail"`
	CustomerName  string `gorm:"type:text;not null" json:"customerName"`

	CustomerID *string `gorm:"type:varchar(50)" json:"customerId,omitempty"`
}
