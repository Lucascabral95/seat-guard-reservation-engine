package services

import (
	"booking-service/internal/models"
	"booking-service/internal/repositories"
)

type CheckoutService struct {
	repo repositories.CheckoutRepository
}

func NewCheckoutService(repo repositories.CheckoutRepository) *CheckoutService {
	return &CheckoutService{repo: repo}
}

func (s *CheckoutService) Create(checkout *models.Checkout) error {
	return s.repo.Create(checkout)
}

func (s *CheckoutService) FindByOrderID(orderID string) (*models.Checkout, error) {
	return s.repo.FindByOrderID(orderID)
}

func (s *CheckoutService) Update(checkout *models.Checkout) error {
	return s.repo.Update(checkout)
}

func (s *CheckoutService) FindAll() ([]models.Checkout, error) {
	return s.repo.FindAll()
}
