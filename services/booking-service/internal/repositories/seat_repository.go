package repositories

import (
	"booking-service/internal/models"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type SeatRepository interface {
	Create(seat *models.Seat) error
	FindAlls() ([]models.Seat, error)
	FindByID(id string) (*models.Seat, error)
	UpdateStatus(id string, status models.SeatStatus) error
	LockSeat(id string, userId string, expiresAt time.Time) error
	UnlockIfExpired(id string, now time.Time) error

	FindSeatByEventId(id string) ([]models.Seat, error)

	FindByIDs(ids []string) ([]models.Seat, error)
}

type seatRepository struct {
	db *gorm.DB
}

func NewSeatRepository(db *gorm.DB) SeatRepository {
	return &seatRepository{db: db}
}

func (r *seatRepository) Create(seat *models.Seat) error {
	return r.db.Create(seat).Error
}

func (r *seatRepository) FindAlls() ([]models.Seat, error) {
	var seats []models.Seat
	err := r.db.Find(&seats).Error
	return seats, err
}

func (r *seatRepository) FindByID(id string) (*models.Seat, error) {
	var seat models.Seat
	err := r.db.First(&seat, "id = ?", id).Error
	return &seat, err
}

func (r *seatRepository) UpdateStatus(id string, status models.SeatStatus) error {
	result := r.db.Model(&models.Seat{}).Where("id = ?", id).Update("status", status)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// Bloquear asiento por 15 minutos
func (r *seatRepository) LockSeat(id, userId string, expiresAt time.Time) error {
	tx := r.db.Model(&models.Seat{}).
		Where("id = ? AND status = ?", id, models.StatusAvailable).
		Updates(map[string]interface{}{
			"status":    models.StatusLocked,
			"locked_by": userId,
			"locked_at": expiresAt,
		})

	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return errors.New("seat not available or not found")
	}
	return nil
}

// Worker que se encarga de verificar si ya paso el tiempo de bloqueo de un asiento
func (r *seatRepository) UnlockIfExpired(id string, now time.Time) error {
	tx := r.db.Model(&models.Seat{}).
		Where("id = ? AND status = ? AND locked_at IS NOT NULL AND locked_at <= ?", id, models.StatusLocked, now).
		Updates(map[string]interface{}{
			"status":    models.StatusAvailable,
			"locked_by": nil,
			"locked_at": nil,
		})

	return tx.Error
}

func (r *seatRepository) FindSeatByEventId(eventId string) ([]models.Seat, error) {
	var seats []models.Seat
	err := r.db.Where("event_id = ?", eventId).Find(&seats).Error
	return seats, err
}

// Muestra los datos de todos los asientos por IDs
func (r *seatRepository) FindByIDs(ids []string) ([]models.Seat, error) {
	if len(ids) == 0 {
		return []models.Seat{}, nil
	}

	var seats []models.Seat

	err := r.db.Where("id IN ? AND deleted_at IS NULL", ids).Find(&seats).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find seats by IDs: %w", err)
	}

	if len(seats) != len(ids) {
		return seats, fmt.Errorf("some seats were not found: expected %d, got %d", len(ids), len(seats))
	}

	return seats, nil
}
