package repositories

import (
	"booking-service/internal/models"
	"errors"

	"gorm.io/gorm"
)

type EventRepository interface {
	Create(event *models.Event) error
	FindByID(id string) (*models.Event, error)
	FindAll() ([]models.Event, error)
	Update(event *models.Event) error
	Delete(id string) error
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

func (r *eventRepository) FindAll() ([]models.Event, error) {
	var events []models.Event
	err := r.db.Find(&events).Error
	return events, err
}

func (r *eventRepository) Update(event *models.Event) error {
	return r.db.Save(event).Error
}

func (r *eventRepository) Delete(id string) error {
	return r.db.Delete(&models.Event{}, "id = ?", id).Error
}
