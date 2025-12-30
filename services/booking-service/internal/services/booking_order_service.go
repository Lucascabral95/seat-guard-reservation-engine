package services

import (
	"booking-service/internal/models"
	"booking-service/internal/repositories"
)

type BookingOrderService struct {
	repo       repositories.BookingOrderRepository
	repoSeats  repositories.SeatRepository
	repoEvents repositories.EventRepository
}

func NewBookingOrderService(repo repositories.BookingOrderRepository, repoSeats repositories.SeatRepository, repoEvents repositories.EventRepository) *BookingOrderService {
	return &BookingOrderService{repo: repo, repoSeats: repoSeats, repoEvents: repoEvents}
}

func (s *BookingOrderService) CreateBookingOrder(bookingOrder *models.BookingOrder) error {
	return s.repo.Create(bookingOrder)
}

func (s *BookingOrderService) FindAllBookingOrders() ([]models.BookingOrder, error) {
	return s.repo.FindAll()
}

func (s *BookingOrderService) FindBookingOrderById(id string) (*models.BookingOrder, error) {
	bookingOrder, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	var detailedSeats []models.Seat

	for _, seatID := range bookingOrder.SeatIDs {
		seat, err := s.repoSeats.FindByID(seatID)
		if err != nil || seat == nil {
			continue
		}

		event, err := s.repoEvents.FindByID(seat.EventID)
		if err == nil && event != nil {
			seat.EventName = event.Name
			seat.EventHour = event.Date.Format("15:04")
		}

		detailedSeats = append(detailedSeats, *seat)
	}

	bookingOrder.Items = detailedSeats

	if len(detailedSeats) > 0 {
		bookingOrder.EventName = detailedSeats[0].EventName
	}

	return bookingOrder, nil
}

func (s *BookingOrderService) UpdateBookingOrderStatus(id string, status models.PaymentStatus) error {
	return s.repo.UpdateStatus(id, status)
}

func (s *BookingOrderService) UpdateBookingOrder(id string, status models.PaymentStatus, paymentProviderID string) error {
	if paymentProviderID == "" {
		return s.repo.UpdateStatus(id, status)
	}
	return s.repo.Update(id, status, paymentProviderID)
}

func (s *BookingOrderService) FindAllOrdersByUserID(userID string) ([]models.BookingOrder, error) {
	bookingOrders, err := s.repo.FindAllOrdersByUserID(userID)
	if err != nil {
		return nil, err
	}

	for i := range bookingOrders {
		var detailedSeats []models.Seat

		for _, seatID := range bookingOrders[i].SeatIDs {
			seat, err := s.repoSeats.FindByID(seatID)
			if err != nil || seat == nil {
				continue
			}

			event, err := s.repoEvents.FindByID(seat.EventID)
			if err == nil && event != nil {
				seat.EventName = event.Name
				seat.EventHour = event.Date.Format("15:04")
			}

			detailedSeats = append(detailedSeats, *seat)
		}

		bookingOrders[i].Items = detailedSeats

		if len(detailedSeats) > 0 {
			bookingOrders[i].EventName = detailedSeats[0].EventName
			bookingOrders[i].EventHour = detailedSeats[0].EventHour
		}
	}

	return bookingOrders, nil
}
