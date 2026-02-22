package services

import (
	"booking-service/internal/models"
	"booking-service/pkg/utils"
	"errors"
	"testing"
	"time"

	"gorm.io/gorm"
)

type mockSeatRepo struct {
	createFn          func(*models.Seat) error
	findAllFn         func() ([]models.Seat, error)
	findByIDFn        func(string) (*models.Seat, error)
	updateStatusFn    func(string, models.SeatStatus) error
	lockSeatFn        func(string, string, time.Time) error
	unlockIfExpiredFn func(string, time.Time) error
	findByEventIDFn   func(string) ([]models.Seat, error)
	findByIDsFn       func([]string) ([]models.Seat, error)
}

func (m *mockSeatRepo) Create(seat *models.Seat) error { return m.createFn(seat) }
func (m *mockSeatRepo) FindAlls() ([]models.Seat, error) { return m.findAllFn() }
func (m *mockSeatRepo) FindByID(id string) (*models.Seat, error) { return m.findByIDFn(id) }
func (m *mockSeatRepo) UpdateStatus(id string, status models.SeatStatus) error {
	return m.updateStatusFn(id, status)
}
func (m *mockSeatRepo) LockSeat(id, userId string, expiresAt time.Time) error {
	return m.lockSeatFn(id, userId, expiresAt)
}
func (m *mockSeatRepo) UnlockIfExpired(id string, now time.Time) error { return m.unlockIfExpiredFn(id, now) }
func (m *mockSeatRepo) FindSeatByEventId(id string) ([]models.Seat, error) { return m.findByEventIDFn(id) }
func (m *mockSeatRepo) FindByIDs(ids []string) ([]models.Seat, error) { return m.findByIDsFn(ids) }

type mockEventRepoForSeat struct {
	findByIDFn func(string) (*models.Event, error)
}

func (m *mockEventRepoForSeat) Create(*models.Event) error { panic("not used") }
func (m *mockEventRepoForSeat) FindByID(id string) (*models.Event, error) { return m.findByIDFn(id) }
func (m *mockEventRepoForSeat) FindAll(models.EventFilter) ([]models.Event, error) { panic("not used") }
func (m *mockEventRepoForSeat) Update(*models.Event) error { panic("not used") }
func (m *mockEventRepoForSeat) Delete(string) error { panic("not used") }
func (m *mockEventRepoForSeat) UpdateAvailability(string) error { panic("not used") }

func TestSeatService_CreateSeat_RejectsNegativePrice(t *testing.T) {
	svc := NewSeatService(&mockSeatRepo{}, &mockEventRepoForSeat{})
	if err := svc.CreateSeat(&models.Seat{Price: -1}); err == nil {
		t.Fatalf("expected validation error")
	}
}

func TestSeatService_GetSeat_MapsRecordNotFound(t *testing.T) {
	svc := NewSeatService(
		&mockSeatRepo{
			unlockIfExpiredFn: func(string, time.Time) error { return nil },
			findByIDFn:        func(string) (*models.Seat, error) { return nil, gorm.ErrRecordNotFound },
		},
		&mockEventRepoForSeat{findByIDFn: func(string) (*models.Event, error) { return nil, nil }},
	)

	_, err := svc.GetSeat("s1")
	if !errors.Is(err, utils.ErrSeatNotFound) {
		t.Fatalf("expected ErrSeatNotFound, got %v", err)
	}
}

func TestSeatService_UpdateSeatStatus_ValidationAndNotFound(t *testing.T) {
	svc := NewSeatService(
		&mockSeatRepo{updateStatusFn: func(string, models.SeatStatus) error { return gorm.ErrRecordNotFound }},
		&mockEventRepoForSeat{},
	)

	if err := svc.UpdateSeatStatus("s1", models.SeatStatus("BAD")); err == nil {
		t.Fatalf("expected invalid status error")
	}

	err := svc.UpdateSeatStatus("s1", models.StatusAvailable)
	if !errors.Is(err, utils.ErrSeatNotFound) {
		t.Fatalf("expected ErrSeatNotFound, got %v", err)
	}
}

func TestSeatService_LockSeat(t *testing.T) {
	t.Run("not available", func(t *testing.T) {
		svc := NewSeatService(
			&mockSeatRepo{
				findByIDFn: func(string) (*models.Seat, error) { return &models.Seat{Status: models.StatusSold}, nil },
			},
			&mockEventRepoForSeat{},
		)
		if err := svc.LockSeat("s1", "u1"); err == nil {
			t.Fatalf("expected error for non-available seat")
		}
	})

	t.Run("available", func(t *testing.T) {
		called := false
		svc := NewSeatService(
			&mockSeatRepo{
				findByIDFn: func(string) (*models.Seat, error) { return &models.Seat{Status: models.StatusAvailable}, nil },
				lockSeatFn: func(id, user string, expires time.Time) error {
					called = true
					if id != "s1" || user != "u1" || time.Until(expires) <= 0 {
						t.Fatalf("unexpected lock args: %s %s %v", id, user, expires)
					}
					return nil
				},
			},
			&mockEventRepoForSeat{},
		)
		if err := svc.LockSeat("s1", "u1"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !called {
			t.Fatalf("expected LockSeat repo call")
		}
	})
}
