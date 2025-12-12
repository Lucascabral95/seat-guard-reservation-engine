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
	return &booking, err
}

func (t *bookingOrderRepository) UpdateStatus(id string, status models.PaymentStatus) error {
	return t.db.Model(&models.BookingOrder{}).Where("id = ?", id).Update("status", status).Error
}
