package services

import (
	"booking-service/internal/models"
	"booking-service/internal/repositories"
)

type BookingOrderService struct {
	repo repositories.BookingOrderRepository
}

func NewBookingOrderService(repo repositories.BookingOrderRepository) *BookingOrderService {
	return &BookingOrderService{repo: repo}
}

func (s *BookingOrderService) CreateBookingOrder(bookingOrder *models.BookingOrder) error {
	return s.repo.Create(bookingOrder)
}

func (s *BookingOrderService) FindAllBookingOrders() ([]models.BookingOrder, error) {
	return s.repo.FindAll()
}

func (s *BookingOrderService) FindBookingOrderById(id string) (*models.BookingOrder, error) {
	return s.repo.FindByID(id)
}

func (s *BookingOrderService) UpdateBookingOrderStatus(id string, status models.PaymentStatus) error {
	return s.repo.UpdateStatus(id, status)
}
