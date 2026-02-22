package services

import (
	"booking-service/internal/models"
	"errors"
	"testing"
	"time"
)

type mockTicketRepo struct {
	createFn        func(*models.TicketPDF) error
	findAllFn       func() ([]*models.TicketPDF, error)
	findByIDFn      func(string) (*models.TicketPDF, error)
	findByOrderIDFn func(string) (*models.TicketPDF, error)
	updateFn        func(*models.TicketPDF) error
	deleteFn        func(string) error
}

func (m *mockTicketRepo) CreateTicket(t *models.TicketPDF) error         { return m.createFn(t) }
func (m *mockTicketRepo) FindAllTickets() ([]*models.TicketPDF, error)   { return m.findAllFn() }
func (m *mockTicketRepo) FindTicketById(id string) (*models.TicketPDF, error) { return m.findByIDFn(id) }
func (m *mockTicketRepo) FindTicketByOrderID(id string) (*models.TicketPDF, error) {
	return m.findByOrderIDFn(id)
}
func (m *mockTicketRepo) UpdateTicket(t *models.TicketPDF) error { return m.updateFn(t) }
func (m *mockTicketRepo) DeleteTicket(id string) error            { return m.deleteFn(id) }

type mockSeatRepoForTicket struct {
	findByIDsFn func([]string) ([]models.Seat, error)
}

func (m *mockSeatRepoForTicket) Create(*models.Seat) error { panic("not used") }
func (m *mockSeatRepoForTicket) FindAlls() ([]models.Seat, error) { panic("not used") }
func (m *mockSeatRepoForTicket) FindByID(string) (*models.Seat, error) { panic("not used") }
func (m *mockSeatRepoForTicket) UpdateStatus(string, models.SeatStatus) error { panic("not used") }
func (m *mockSeatRepoForTicket) LockSeat(string, string, time.Time) error { panic("not used") }
func (m *mockSeatRepoForTicket) UnlockIfExpired(string, time.Time) error { panic("not used") }
func (m *mockSeatRepoForTicket) FindSeatByEventId(string) ([]models.Seat, error) { panic("not used") }
func (m *mockSeatRepoForTicket) FindByIDs(ids []string) ([]models.Seat, error) { return m.findByIDsFn(ids) }

type mockOrderRepoForTicket struct {
	findByIDFn func(string) (*models.BookingOrder, error)
}

func (m *mockOrderRepoForTicket) Create(*models.BookingOrder) error { panic("not used") }
func (m *mockOrderRepoForTicket) FindAll() ([]models.BookingOrder, error) { panic("not used") }
func (m *mockOrderRepoForTicket) FindByID(id string) (*models.BookingOrder, error) { return m.findByIDFn(id) }
func (m *mockOrderRepoForTicket) UpdateStatus(string, models.PaymentStatus) error { panic("not used") }
func (m *mockOrderRepoForTicket) Update(string, models.PaymentStatus, string) error { panic("not used") }
func (m *mockOrderRepoForTicket) FindAllOrdersByUserID(string) ([]models.BookingOrder, error) { panic("not used") }

type mockEventRepoForTicket struct {
	findByIDFn func(string) (*models.Event, error)
}

func (m *mockEventRepoForTicket) Create(*models.Event) error { panic("not used") }
func (m *mockEventRepoForTicket) FindByID(id string) (*models.Event, error) { return m.findByIDFn(id) }
func (m *mockEventRepoForTicket) FindAll(models.EventFilter) ([]models.Event, error) { panic("not used") }
func (m *mockEventRepoForTicket) Update(*models.Event) error { panic("not used") }
func (m *mockEventRepoForTicket) Delete(string) error { panic("not used") }
func (m *mockEventRepoForTicket) UpdateAvailability(string) error { panic("not used") }

func TestTicketService_CreateTicketFromOrder_ValidationsAndSuccess(t *testing.T) {
	svc := NewTicketService(&mockTicketRepo{}, &mockOrderRepoForTicket{}, &mockSeatRepoForTicket{}, &mockEventRepoForTicket{})
	if _, err := svc.CreateTicketFromOrder(nil, &models.BookingOrder{}); err == nil {
		t.Fatalf("expected nil checkout/order validation error")
	}

	createCalled := false
	svc = NewTicketService(
		&mockTicketRepo{createFn: func(ticket *models.TicketPDF) error {
			createCalled = true
			if ticket.EventName != "Rock Fest" || ticket.EventHour != "19:45" {
				t.Fatalf("unexpected event mapping: %+v", ticket)
			}
			return nil
		}},
		&mockOrderRepoForTicket{},
		&mockSeatRepoForTicket{findByIDsFn: func([]string) ([]models.Seat, error) {
			return []models.Seat{{EventID: "e1", Number: "A1"}}, nil
		}},
		&mockEventRepoForTicket{findByIDFn: func(string) (*models.Event, error) {
			return &models.Event{Name: "Rock Fest", Date: time.Date(2026, 6, 1, 19, 45, 0, 0, time.UTC)}, nil
		}},
	)

	checkout := &models.Checkout{PaymentProvider: "STRIPE", PaymentIntentID: "pi_1", Currency: "USD", Amount: 1000, CustomerName: "Ana", CustomerEmail: "a@a.com"}
	order := &models.BookingOrder{BaseModel: models.BaseModel{ID: "o1"}, SeatIDs: []string{"s1"}}
	got, err := svc.CreateTicketFromOrder(checkout, order)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !createCalled || got.OrderID != "o1" {
		t.Fatalf("expected ticket creation and order mapping")
	}
}

func TestTicketService_UpdateTicketPDF_IncrementsVersion(t *testing.T) {
	updated := &models.TicketPDF{}
	svc := NewTicketService(
		&mockTicketRepo{
			findByIDFn: func(string) (*models.TicketPDF, error) { return &models.TicketPDF{BaseModel: models.BaseModel{ID: "t1"}, PDFVersion: 1}, nil },
			updateFn:   func(ticket *models.TicketPDF) error { *updated = *ticket; return nil },
		},
		&mockOrderRepoForTicket{},
		&mockSeatRepoForTicket{},
		&mockEventRepoForTicket{},
	)

	if err := svc.UpdateTicketPDF("t1", []byte("pdf")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.PDFVersion != 2 || string(updated.PDFData) != "pdf" || updated.PDFGeneratedAt == nil {
		t.Fatalf("expected ticket update with incremented version, got %+v", updated)
	}
}

func TestTicketService_ValidateTicketOwnership(t *testing.T) {
	svc := NewTicketService(
		&mockTicketRepo{findByIDFn: func(string) (*models.TicketPDF, error) { return &models.TicketPDF{OrderID: "o1"}, nil }},
		&mockOrderRepoForTicket{findByIDFn: func(string) (*models.BookingOrder, error) { return &models.BookingOrder{UserID: "u2"}, nil }},
		&mockSeatRepoForTicket{},
		&mockEventRepoForTicket{},
	)

	if err := svc.ValidateTicketOwnership("t1", "u1"); err == nil {
		t.Fatalf("expected access denied")
	}
}

func TestTicketService_GetTicketByID_LoadsItems(t *testing.T) {
	svc := NewTicketService(
		&mockTicketRepo{findByIDFn: func(string) (*models.TicketPDF, error) { return &models.TicketPDF{OrderID: "o1"}, nil }},
		&mockOrderRepoForTicket{findByIDFn: func(string) (*models.BookingOrder, error) { return &models.BookingOrder{SeatIDs: []string{"s1"}}, nil }},
		&mockSeatRepoForTicket{findByIDsFn: func([]string) ([]models.Seat, error) { return []models.Seat{{EventID: "e1", Number: "A1"}}, nil }},
		&mockEventRepoForTicket{findByIDFn: func(string) (*models.Event, error) {
			return &models.Event{Name: "Show", Date: time.Date(2026, 2, 1, 18, 0, 0, 0, time.UTC)}, nil
		}},
	)

	ticket, err := svc.GetTicketByID("t1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ticket.Items) != 1 || ticket.EventName != "Show" || ticket.EventHour != "18:00" {
		t.Fatalf("expected loaded items/event fields, got %+v", ticket)
	}
}

func TestTicketService_CreateTicketFromOrder_SeatFetchFailure(t *testing.T) {
	svc := NewTicketService(
		&mockTicketRepo{},
		&mockOrderRepoForTicket{},
		&mockSeatRepoForTicket{findByIDsFn: func([]string) ([]models.Seat, error) { return nil, errors.New("db") }},
		&mockEventRepoForTicket{},
	)

	_, err := svc.CreateTicketFromOrder(&models.Checkout{}, &models.BookingOrder{SeatIDs: []string{"s1"}})
	if err == nil {
		t.Fatalf("expected seat fetch failure")
	}
}
