package services

import (
	"booking-service/internal/models"
	"booking-service/internal/repositories"
	"errors"
	"time"
)

type SeatService struct {
	repo repositories.SeatRepository
}

func NewSeatService(repo repositories.SeatRepository) *SeatService {
	return &SeatService{repo: repo}
}

func (s *SeatService) CreateSeat(seat *models.Seat) error {
	if seat.Price < 0 {
		return errors.New("seat price cannot be negative")
	}

	return s.repo.Create(seat)
}

func (s *SeatService) GetSeats() ([]models.Seat, error) {
	return s.repo.FindAlls()
}

func (s *SeatService) GetSeat(id string) (*models.Seat, error) {
	if err := s.repo.UnlockIfExpired(id, time.Now()); err != nil {
		return nil, err
	}

	seat, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return seat, nil
}

func (s *SeatService) UpdateSeatStatus(id string, status models.SeatStatus) error {
	if status != models.StatusAvailable && status != models.StatusLocked && status != models.StatusSold {
		return errors.New("invalid seat status")
	}

	return s.repo.UpdateStatus(id, status)
}

func (s *SeatService) GetSeatByEventId(eventId string) ([]models.Seat, error) {
	return s.repo.FindSeatByEventId(eventId)
}

// Bloquear asiento por 15 minutos
func (s *SeatService) LockSeat(id string, userId string) error {
	seat, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("seat not found")
	}

	if seat.Status != models.StatusAvailable {
		return errors.New("seat is not available")
	}

	// Bloquear por 15 minutos
	expiresAt := time.Now().Add(15 * time.Minute)
	return s.repo.LockSeat(id, userId, expiresAt)
}
