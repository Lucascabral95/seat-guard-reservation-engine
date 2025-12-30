package repositories

import (
	"booking-service/internal/models"
	"errors"

	"gorm.io/gorm"
)

type EventRepository interface {
	Create(event *models.Event) error
	FindByID(id string) (*models.Event, error)
	FindAll(filter models.EventFilter) ([]models.Event, error)
	Update(event *models.Event) error
	Delete(id string) error

	UpdateAvailability(eventID string) error
}

type eventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) EventRepository {
	return &eventRepository{db: db}
}

func (r *eventRepository) Create(event *models.Event) error {
	return r.db.Create(event).Error
}

func (r *eventRepository) FindByID(id string) (*models.Event, error) {
	var event models.Event
	err := r.db.Preload("Seats").First(&event, "id = ?", id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &event, err
}

func (r *eventRepository) FindAll(filter models.EventFilter) ([]models.Event, error) {
	var events []models.Event

	query := r.db.Model(&models.Event{})

	if filter.Name != "" {
		query = query.Where("name ILIKE ?", "%"+filter.Name+"%")
	}

	if filter.Gender != "" {
		query = query.Where("gender = ?", filter.Gender)
	}

	if filter.Location != "" {
		query = query.Where("location ILIKE ?", "%"+filter.Location+"%")
	}

	err := query.Find(&events).Error
	return events, err
}

func (r *eventRepository) Update(event *models.Event) error {
	return r.db.Save(event).Error
}

func (r *eventRepository) Delete(id string) error {
	return r.db.Delete(&models.Event{}, "id = ?", id).Error
}

func (r *eventRepository) UpdateAvailability(eventID string) error {
	var totalSeats int64
	var availableSeats int64

	if err := r.db.Model(&models.Seat{}).Where("event_id = ?", eventID).Count(&totalSeats).Error; err != nil {
		return err
	}

	if err := r.db.Model(&models.Seat{}).Where("event_id", eventID).Count(&availableSeats).Error; err != nil {
		return err
	}

	if totalSeats == 0 {
		return nil
	}
	percentage := float64(availableSeats) / float64(totalSeats)
	var newStatus models.Availability

	switch {
	case availableSeats == 0:
		newStatus = models.AvailabilitySoldOut
	case percentage > 0.5:
		newStatus = models.AvailabilityHigh
	case availableSeats < 5:
		newStatus = models.AvailabilityLow
	default:
		newStatus = models.AvailabilityMedium
	}

	return r.db.Model(&models.Event{}).Where("id = ?", eventID).Update("availability", newStatus).Error
}
