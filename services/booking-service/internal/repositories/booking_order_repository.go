package repositories

import (
	"booking-service/internal/models"

	"gorm.io/gorm"
)

type BookingOrderRepository interface {
	Create(booking *models.BookingOrder) error
	FindAll() ([]models.BookingOrder, error)
	FindByID(id string) (*models.BookingOrder, error)
	UpdateStatus(id string, status models.PaymentStatus) error
	Update(id string, status models.PaymentStatus, paymentProviderID string) error

	FindAllOrdersByUserID(userID string) ([]models.BookingOrder, error)
}

type bookingOrderRepository struct {
	db *gorm.DB
}

func NewBookingOrderRepository(db *gorm.DB) BookingOrderRepository {
	return &bookingOrderRepository{db: db}
}

func (t *bookingOrderRepository) Create(bookingOrder *models.BookingOrder) error {
	return t.db.Create(bookingOrder).Error
}

func (t *bookingOrderRepository) FindAll() ([]models.BookingOrder, error) {
	var bookings []models.BookingOrder
	err := t.db.Find(&bookings).Error
	return bookings, err
}

func (t *bookingOrderRepository) FindByID(id string) (*models.BookingOrder, error) {
	var booking models.BookingOrder
	err := t.db.First(&booking, "id = ?", id).Error
	if err != nil {
		return &booking, err
	}

	if len(booking.SeatIDs) == 0 {
		return &booking, nil
	}

	var seats []models.Seat
	err = t.db.Where("id IN ?", booking.SeatIDs).Find(&seats).Error
	if err != nil {
		return &booking, err
	}

	seatByID := make(map[string]models.Seat, len(seats))
	for _, s := range seats {
		seatByID[s.ID] = s
	}

	booking.Items = booking.Items[:0]
	for _, seatID := range booking.SeatIDs {
		if seat, ok := seatByID[seatID]; ok {
			booking.Items = append(booking.Items, seat)
		}
	}

	return &booking, nil
}

func (t *bookingOrderRepository) UpdateStatus(id string, status models.PaymentStatus) error {
	return t.db.Model(&models.BookingOrder{}).Where("id = ?", id).Update("status", status).Error
}

func (t *bookingOrderRepository) Update(id string, status models.PaymentStatus, paymentProviderID string) error {
	updates := map[string]any{
		"status":              status,
		"payment_provider_id": paymentProviderID,
	}
	return t.db.Model(&models.BookingOrder{}).Where("id = ?", id).Updates(updates).Error
}

func (t *bookingOrderRepository) FindAllOrdersByUserID(userID string) ([]models.BookingOrder, error) {
	var bookings []models.BookingOrder
	err := t.db.Where("user_id = ?", userID).Find(&bookings).Error
	return bookings, err
}
