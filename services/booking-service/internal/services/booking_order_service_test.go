package services

import (
	"booking-service/internal/models"
	"errors"
	"testing"
	"time"
)

type mockBookingOrderRepo struct {
	createFn                func(*models.BookingOrder) error
	findAllFn               func() ([]models.BookingOrder, error)
	findByIDFn              func(string) (*models.BookingOrder, error)
	updateStatusFn          func(string, models.PaymentStatus) error
	updateFn                func(string, models.PaymentStatus, string) error
	findAllOrdersByUserIDFn func(string) ([]models.BookingOrder, error)
}

func (m *mockBookingOrderRepo) Create(booking *models.BookingOrder) error { return m.createFn(booking) }
func (m *mockBookingOrderRepo) FindAll() ([]models.BookingOrder, error)   { return m.findAllFn() }
func (m *mockBookingOrderRepo) FindByID(id string) (*models.BookingOrder, error) {
	return m.findByIDFn(id)
}
func (m *mockBookingOrderRepo) UpdateStatus(id string, status models.PaymentStatus) error {
	return m.updateStatusFn(id, status)
}
func (m *mockBookingOrderRepo) Update(id string, status models.PaymentStatus, paymentProviderID string) error {
	return m.updateFn(id, status, paymentProviderID)
}
func (m *mockBookingOrderRepo) FindAllOrdersByUserID(userID string) ([]models.BookingOrder, error) {
	return m.findAllOrdersByUserIDFn(userID)
}

type mockSeatRepoForBooking struct {
	findByIDFn func(string) (*models.Seat, error)
}

func (m *mockSeatRepoForBooking) Create(*models.Seat) error { panic("not used") }
func (m *mockSeatRepoForBooking) FindAlls() ([]models.Seat, error) { panic("not used") }
func (m *mockSeatRepoForBooking) FindByID(id string) (*models.Seat, error) { return m.findByIDFn(id) }
func (m *mockSeatRepoForBooking) UpdateStatus(string, models.SeatStatus) error { panic("not used") }
func (m *mockSeatRepoForBooking) LockSeat(string, string, time.Time) error { panic("not used") }
func (m *mockSeatRepoForBooking) UnlockIfExpired(string, time.Time) error { panic("not used") }
func (m *mockSeatRepoForBooking) FindSeatByEventId(string) ([]models.Seat, error) { panic("not used") }
func (m *mockSeatRepoForBooking) FindByIDs([]string) ([]models.Seat, error) { panic("not used") }

type mockEventRepoForBooking struct {
	findByIDFn func(string) (*models.Event, error)
}

func (m *mockEventRepoForBooking) Create(*models.Event) error { panic("not used") }
func (m *mockEventRepoForBooking) FindByID(id string) (*models.Event, error) { return m.findByIDFn(id) }
func (m *mockEventRepoForBooking) FindAll(models.EventFilter) ([]models.Event, error) { panic("not used") }
func (m *mockEventRepoForBooking) Update(*models.Event) error { panic("not used") }
func (m *mockEventRepoForBooking) Delete(string) error { panic("not used") }
func (m *mockEventRepoForBooking) UpdateAvailability(string) error { panic("not used") }

func TestBookingOrderService_UpdateBookingOrder_UsesCorrectRepoMethod(t *testing.T) {
	statusCalled := 0
	updateCalled := 0

	svc := NewBookingOrderService(
		&mockBookingOrderRepo{
			updateStatusFn: func(string, models.PaymentStatus) error {
				statusCalled++
				return nil
			},
			updateFn: func(string, models.PaymentStatus, string) error {
				updateCalled++
				return nil
			},
		},
		&mockSeatRepoForBooking{},
		&mockEventRepoForBooking{},
	)

	if err := svc.UpdateBookingOrder("o1", models.PaymentCompleted, ""); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if statusCalled != 1 || updateCalled != 0 {
		t.Fatalf("expected UpdateStatus once, got status=%d update=%d", statusCalled, updateCalled)
	}

	if err := svc.UpdateBookingOrder("o1", models.PaymentCompleted, "pi_123"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if statusCalled != 1 || updateCalled != 1 {
		t.Fatalf("expected Update once on second call, got status=%d update=%d", statusCalled, updateCalled)
	}
}

func TestBookingOrderService_FindBookingOrderById_EnrichesSeatsAndEvent(t *testing.T) {
	service := NewBookingOrderService(
		&mockBookingOrderRepo{
			findByIDFn: func(id string) (*models.BookingOrder, error) {
				return &models.BookingOrder{BaseModel: models.BaseModel{ID: id}, SeatIDs: []string{"s1", "missing"}}, nil
			},
		},
		&mockSeatRepoForBooking{findByIDFn: func(id string) (*models.Seat, error) {
			if id == "missing" {
				return nil, errors.New("not found")
			}
			return &models.Seat{BaseModel: models.BaseModel{ID: "s1"}, EventID: "e1"}, nil
		}},
		&mockEventRepoForBooking{findByIDFn: func(string) (*models.Event, error) {
			return &models.Event{Name: "Show", Date: time.Date(2026, 3, 1, 21, 30, 0, 0, time.UTC)}, nil
		}},
	)

	got, err := service.FindBookingOrderById("o1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got.Items) != 1 {
		t.Fatalf("expected 1 valid seat, got %d", len(got.Items))
	}
	if got.Items[0].EventName != "Show" || got.Items[0].EventHour != "21:30" {
		t.Fatalf("expected enriched seat event data, got %#v", got.Items[0])
	}
	if got.EventName != "Show" {
		t.Fatalf("expected booking order event name to be set")
	}
}

func TestBookingOrderService_FindAllOrdersByUserID_EnrichesOrders(t *testing.T) {
	svc := NewBookingOrderService(
		&mockBookingOrderRepo{findAllOrdersByUserIDFn: func(string) ([]models.BookingOrder, error) {
			return []models.BookingOrder{{BaseModel: models.BaseModel{ID: "o1"}, SeatIDs: []string{"s1"}}}, nil
		}},
		&mockSeatRepoForBooking{findByIDFn: func(string) (*models.Seat, error) {
			return &models.Seat{BaseModel: models.BaseModel{ID: "s1"}, EventID: "e1"}, nil
		}},
		&mockEventRepoForBooking{findByIDFn: func(string) (*models.Event, error) {
			return &models.Event{Name: "Concert", Date: time.Date(2026, 5, 10, 20, 0, 0, 0, time.UTC)}, nil
		}},
	)

	orders, err := svc.FindAllOrdersByUserID("u1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(orders) != 1 || len(orders[0].Items) != 1 {
		t.Fatalf("expected one order with one enriched item")
	}
	if orders[0].EventName != "Concert" || orders[0].EventHour != "20:00" {
		t.Fatalf("expected order event fields enriched, got %+v", orders[0])
	}
}
