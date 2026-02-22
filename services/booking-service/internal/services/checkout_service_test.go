package services

import (
	"booking-service/internal/models"
	"errors"
	"testing"
)

type mockCheckoutRepo struct {
	createFn      func(*models.Checkout) error
	findByOrderFn func(string) (*models.Checkout, error)
	updateFn      func(*models.Checkout) error
	findAllFn     func() ([]models.Checkout, error)
}

func (m *mockCheckoutRepo) Create(c *models.Checkout) error         { return m.createFn(c) }
func (m *mockCheckoutRepo) FindByOrderID(id string) (*models.Checkout, error) { return m.findByOrderFn(id) }
func (m *mockCheckoutRepo) Update(c *models.Checkout) error         { return m.updateFn(c) }
func (m *mockCheckoutRepo) FindAll() ([]models.Checkout, error)     { return m.findAllFn() }

func TestCheckoutService_PassThroughMethods(t *testing.T) {
	expected := errors.New("boom")
	svc := NewCheckoutService(&mockCheckoutRepo{
		createFn: func(*models.Checkout) error { return expected },
		findByOrderFn: func(id string) (*models.Checkout, error) {
			return &models.Checkout{OrderID: id}, nil
		},
		updateFn:  func(*models.Checkout) error { return nil },
		findAllFn: func() ([]models.Checkout, error) { return []models.Checkout{{OrderID: "o1"}}, nil },
	})

	if err := svc.Create(&models.Checkout{}); !errors.Is(err, expected) {
		t.Fatalf("expected repo error, got %v", err)
	}

	c, err := svc.FindByOrderID("o1")
	if err != nil || c.OrderID != "o1" {
		t.Fatalf("unexpected findByOrderID result: %+v err=%v", c, err)
	}

	if err := svc.Update(&models.Checkout{}); err != nil {
		t.Fatalf("unexpected update error: %v", err)
	}

	all, err := svc.FindAll()
	if err != nil || len(all) != 1 {
		t.Fatalf("unexpected findAll result len=%d err=%v", len(all), err)
	}
}
