package services

import (
	"booking-service/internal/models"
	"testing"
	"time"
)

type mockEventRepo struct {
	createFn             func(*models.Event) error
	findByIDFn           func(string) (*models.Event, error)
	findAllFn            func(models.EventFilter) ([]models.Event, error)
	updateFn             func(*models.Event) error
	deleteFn             func(string) error
	updateAvailabilityFn func(string) error
}

func (m *mockEventRepo) Create(e *models.Event) error                               { return m.createFn(e) }
func (m *mockEventRepo) FindByID(id string) (*models.Event, error)                  { return m.findByIDFn(id) }
func (m *mockEventRepo) FindAll(f models.EventFilter) ([]models.Event, error)       { return m.findAllFn(f) }
func (m *mockEventRepo) Update(e *models.Event) error                               { return m.updateFn(e) }
func (m *mockEventRepo) Delete(id string) error                                      { return m.deleteFn(id) }
func (m *mockEventRepo) UpdateAvailability(eventID string) error                     { return m.updateAvailabilityFn(eventID) }

func TestEventService_CreateEvent_Validations(t *testing.T) {
	svc := NewEventService(&mockEventRepo{})

	cases := []models.Event{
		{Name: "", Price: 10, Date: time.Now().Add(time.Hour)},
		{Name: "x", Price: -1, Date: time.Now().Add(time.Hour)},
		{Name: "x", Price: 1, Date: time.Now().Add(-time.Hour)},
	}
	for _, c := range cases {
		if err := svc.CreateEvent(&c); err == nil {
			t.Fatalf("expected validation error for %+v", c)
		}
	}
}

func TestEventService_UpdateEvent_NotFoundAndSuccess(t *testing.T) {
	t.Run("not found", func(t *testing.T) {
		svc := NewEventService(&mockEventRepo{findByIDFn: func(string) (*models.Event, error) { return nil, nil }})
		if err := svc.UpdateEvent("e1", &models.Event{}); err == nil {
			t.Fatalf("expected not found error")
		}
	})

	t.Run("success", func(t *testing.T) {
		updated := &models.Event{}
		svc := NewEventService(&mockEventRepo{
			findByIDFn: func(string) (*models.Event, error) {
				return &models.Event{Name: "old", Price: 1}, nil
			},
			updateFn: func(e *models.Event) error { *updated = *e; return nil },
		})

		err := svc.UpdateEvent("e1", &models.Event{Name: "new", Description: "d", Location: "loc", Date: time.Now().Add(time.Hour), Price: 2})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if updated.Name != "new" || updated.Location != "loc" || updated.Price != 2 {
			t.Fatalf("expected updated fields copied, got %+v", updated)
		}
	})
}
