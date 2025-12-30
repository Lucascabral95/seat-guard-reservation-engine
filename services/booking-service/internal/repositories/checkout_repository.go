package repositories

import (
	"booking-service/internal/models"

	"gorm.io/gorm"
)

type CheckoutRepository interface {
	Create(checkout *models.Checkout) error
	FindByOrderID(orderID string) (*models.Checkout, error)
	Update(checkout *models.Checkout) error
	FindAll() ([]models.Checkout, error)
}

type checkoutRepository struct {
	db *gorm.DB
}

func NewCheckoutRepository(db *gorm.DB) CheckoutRepository {
	return &checkoutRepository{db: db}
}

func (r *checkoutRepository) Create(checkout *models.Checkout) error {
	return r.db.Create(checkout).Error
}

func (r *checkoutRepository) FindByOrderID(orderID string) (*models.Checkout, error) {
	var checkout models.Checkout

	err := r.db.
		Preload("Order").
		Where("order_id = ?", orderID).
		First(&checkout).Error
	if err != nil {
		return nil, err
	}

	var seats []models.Seat
	if len(checkout.Order.SeatIDs) > 0 {
		err = r.db.Where("id IN ?", checkout.Order.SeatIDs).Find(&seats).Error
		if err != nil {
			return nil, err
		}
		checkout.Order.Items = seats
	}

	return &checkout, nil
}

func (r *checkoutRepository) Update(checkout *models.Checkout) error {
	return r.db.Save(checkout).Error
}

func (r *checkoutRepository) FindByID(id string) (*models.Checkout, error) {
	var checkout models.Checkout
	err := r.db.Where("id = ?", id).First(&checkout).Error
	return &checkout, err
}

func (r *checkoutRepository) FindAll() ([]models.Checkout, error) {
	var checkouts []models.Checkout
	err := r.db.Find(&checkouts).Error
	return checkouts, err
}
