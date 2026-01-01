package services

import (
	"booking-service/internal/models"
	"booking-service/internal/repositories"
	"errors"
	"fmt"
	"time"
)

type TicketService struct {
	ticketRepo repositories.TicketRepository
	orderRepo  repositories.BookingOrderRepository
	seatRepo   repositories.SeatRepository
	eventRepo  repositories.EventRepository
}

func NewTicketService(
	ticketRepo repositories.TicketRepository,
	orderRepo repositories.BookingOrderRepository,
	seatRepo repositories.SeatRepository,
	eventRepo repositories.EventRepository,
) *TicketService {
	return &TicketService{
		ticketRepo: ticketRepo,
		orderRepo:  orderRepo,
		seatRepo:   seatRepo,
		eventRepo:  eventRepo,
	}
}

// CreateTicketFromOrder crea un TicketPDF a partir de un Checkout y BookingOrder
func (s *TicketService) CreateTicketFromOrder(
	checkout *models.Checkout,
	order *models.BookingOrder,
) (*models.TicketPDF, error) {
	if checkout == nil || order == nil {
		return nil, errors.New("checkout and order are required")
	}

	seats, err := s.seatRepo.FindByIDs(order.SeatIDs)
	if err != nil || len(seats) == 0 {
		return nil, fmt.Errorf("failed to fetch seats: %w", err)
	}

	eventID := seats[0].EventID
	event, err := s.eventRepo.FindByID(eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch event: %w", err)
	}

	ticket := &models.TicketPDF{
		PaymentProvider: checkout.PaymentProvider,
		PaymentIntentID: checkout.PaymentIntentID,
		Currency:        checkout.Currency,
		Amount:          checkout.Amount,
		Name:            checkout.CustomerName,
		Email:           checkout.CustomerEmail,
		CustomerID:      checkout.CustomerID,
		OrderID:         order.ID,
		EventName:       event.Name,
		EventHour:       event.Date.Format("15:04"),
		Items:           seats,
		PDFVersion:      1,
	}

	if err := s.ticketRepo.CreateTicket(ticket); err != nil {
		return nil, fmt.Errorf("failed to create ticket: %w", err)
	}

	return ticket, nil
}

// GetTicketByID obtiene un ticket por ID y carga sus items
func (s *TicketService) GetTicketByID(ticketID string) (*models.TicketPDF, error) {
	ticket, err := s.ticketRepo.FindTicketById(ticketID)
	if err != nil {
		return nil, err
	}

	if err := s.loadTicketItems(ticket); err != nil {
		return nil, err
	}

	return ticket, nil
}

// GetTicketByOrderID obtiene un ticket por OrderID y carga sus items
func (s *TicketService) GetTicketByOrderID(orderID string) (*models.TicketPDF, error) {
	ticket, err := s.ticketRepo.FindTicketByOrderID(orderID)
	if err != nil {
		return nil, err
	}

	if err := s.loadTicketItems(ticket); err != nil {
		return nil, err
	}

	return ticket, nil
}

// GetAllTickets obtiene todos los tickets
func (s *TicketService) GetAllTickets() ([]*models.TicketPDF, error) {
	return s.ticketRepo.FindAllTickets()
}

// UpdateTicketPDF actualiza el PDF binario de un ticket
func (s *TicketService) UpdateTicketPDF(ticketID string, pdfData []byte) error {
	ticket, err := s.ticketRepo.FindTicketById(ticketID)
	if err != nil {
		return err
	}

	now := time.Now()
	ticket.PDFData = pdfData
	ticket.PDFGeneratedAt = &now
	ticket.PDFVersion++

	return s.ticketRepo.UpdateTicket(ticket)
}

// DeleteTicket elimina un ticket (soft delete)
func (s *TicketService) DeleteTicket(ticketID string) error {
	return s.ticketRepo.DeleteTicket(ticketID)
}

// ValidateTicketOwnership verifica que el ticket pertenece al usuario
func (s *TicketService) ValidateTicketOwnership(ticketID, userID string) error {
	ticket, err := s.ticketRepo.FindTicketById(ticketID)
	if err != nil {
		return err
	}

	order, err := s.orderRepo.FindByID(ticket.OrderID)
	if err != nil {
		return fmt.Errorf("failed to fetch order: %w", err)
	}

	if order.UserID != userID {
		return errors.New("access denied: ticket does not belong to user")
	}

	return nil
}

// loadTicketItems carga los asientos y datos del evento para un ticket
func (s *TicketService) loadTicketItems(ticket *models.TicketPDF) error {
	order, err := s.orderRepo.FindByID(ticket.OrderID)
	if err != nil {
		return fmt.Errorf("failed to fetch order: %w", err)
	}

	seats, err := s.seatRepo.FindByIDs(order.SeatIDs)
	if err != nil {
		return fmt.Errorf("failed to fetch seats: %w", err)
	}

	if len(seats) > 0 {
		event, err := s.eventRepo.FindByID(seats[0].EventID)
		if err != nil {
			return fmt.Errorf("failed to fetch event: %w", err)
		}

		ticket.EventName = event.Name
		ticket.EventHour = event.Date.Format("15:04")
	}

	ticket.Items = seats
	return nil
}
