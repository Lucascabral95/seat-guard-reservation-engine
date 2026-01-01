package repositories

import (
	"booking-service/internal/models"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// TicketRepository define el contrato para operaciones de TicketPDF
type TicketRepository interface {
	CreateTicket(ticket *models.TicketPDF) error

	FindAllTickets() ([]*models.TicketPDF, error)
	FindTicketById(id string) (*models.TicketPDF, error)
	FindTicketByOrderID(orderID string) (*models.TicketPDF, error)

	UpdateTicket(ticket *models.TicketPDF) error
	DeleteTicket(id string) error
}

// ticketRepository es la implementación concreta
type ticketRepository struct {
	db *gorm.DB
}

// NewTicketRepository crea una nueva instancia del repositorio
func NewTicketRepository(db *gorm.DB) TicketRepository {
	return &ticketRepository{db: db}
}

// CreateTicket crea un nuevo ticket en la base de datos
func (r *ticketRepository) CreateTicket(ticket *models.TicketPDF) error {
	if ticket == nil {
		return errors.New("ticket cannot be nil")
	}

	if err := r.db.Create(ticket).Error; err != nil {
		return fmt.Errorf("failed to create ticket: %w", err)
	}

	return nil
}

// FindAllTickets obtiene todos los tickets (sin el PDF binario por defecto)
func (r *ticketRepository) FindAllTickets() ([]*models.TicketPDF, error) {
	var tickets []*models.TicketPDF

	err := r.db.
		Select("id, payment_provider, payment_intent_id, currency, amount, name, email, customer_id, order_id, pdf_generated_at, pdf_version, created_at, updated_at").
		Where("deleted_at IS NULL").
		Find(&tickets).Error

	if err != nil {
		return nil, fmt.Errorf("failed to find tickets: %w", err)
	}

	return tickets, nil
}

// FindTicketById obtiene un ticket por ID (con PDF binario)
func (r *ticketRepository) FindTicketById(id string) (*models.TicketPDF, error) {
	if id == "" {
		return nil, errors.New("ticket ID cannot be empty")
	}

	var ticket models.TicketPDF

	err := r.db.First(&ticket, "id = ? AND deleted_at IS NULL", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("ticket with ID %s not found", id)
		}
		return nil, fmt.Errorf("failed to find ticket: %w", err)
	}

	return &ticket, nil
}

// FindTicketByOrderID obtiene un ticket por OrderID (con PDF binario)
func (r *ticketRepository) FindTicketByOrderID(orderID string) (*models.TicketPDF, error) {
	if orderID == "" {
		return nil, errors.New("order ID cannot be empty")
	}

	var ticket models.TicketPDF

	err := r.db.First(&ticket, "order_id = ? AND deleted_at IS NULL", orderID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("ticket for order %s not found", orderID)
		}
		return nil, fmt.Errorf("failed to find ticket: %w", err)
	}

	return &ticket, nil
}

// UpdateTicket actualiza un ticket existente
func (r *ticketRepository) UpdateTicket(ticket *models.TicketPDF) error {
	if ticket == nil {
		return errors.New("ticket cannot be nil")
	}

	if ticket.ID == "" {
		return errors.New("ticket ID is required for update")
	}

	var existingTicket models.TicketPDF
	if err := r.db.First(&existingTicket, "id = ? AND deleted_at IS NULL", ticket.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("ticket with ID %s not found", ticket.ID)
		}
		return fmt.Errorf("failed to check ticket existence: %w", err)
	}

	if err := r.db.Save(ticket).Error; err != nil {
		return fmt.Errorf("failed to update ticket: %w", err)
	}

	return nil
}

// DeleteTicket realiza soft delete de un ticket
func (r *ticketRepository) DeleteTicket(id string) error {
	if id == "" {
		return errors.New("ticket ID cannot be empty")
	}

	result := r.db.Delete(&models.TicketPDF{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete ticket: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("ticket with ID %s not found", id)
	}

	return nil
}
