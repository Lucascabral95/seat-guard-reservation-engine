package services

import (
	"booking-service/internal/models"
	"booking-service/internal/repositories"
	"errors"
	"time"
)

type EventService struct {
	repo repositories.EventRepository
}

func NewEventService(repo repositories.EventRepository) *EventService {
	return &EventService{repo: repo}
}

func (s *EventService) CreateEvent(event *models.Event) error {
	if event.Name == "" {
		return errors.New("event name cannot be empty")
	}
	if event.Price < 0 {
		return errors.New("event price cannot be negative")
	}
	if event.Date.Before(time.Now()) {
		return errors.New("event date cannot be in the past")
	}

	return s.repo.Create(event)
}

func (s *EventService) GetEvent(id string) (*models.Event, error) {
	event, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if event == nil {
		return nil, errors.New("event not found")
	}

	return event, nil
}

func (s *EventService) GetAllEvents() ([]models.Event, error) {
	return s.repo.FindAll()
}

func (s *EventService) UpdateEvent(id string, updatedData *models.Event) error {
	existingEvent, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if existingEvent == nil {
		return errors.New("cannot update: event not found")
	}

	existingEvent.Name = updatedData.Name
	existingEvent.Description = updatedData.Description
	existingEvent.Location = updatedData.Location
	existingEvent.Date = updatedData.Date
	existingEvent.Price = updatedData.Price

	return s.repo.Update(existingEvent)
}

func (s *EventService) DeleteEvent(id string) error {
	return s.repo.Delete(id)
}
