package services

import (
	"booking-service/internal/models"

	"gorm.io/gorm"
)

func UpdateEventAvailability(db *gorm.DB, eventID string) error {
	var totalSeats int64
	var availableSeats int64

 if eventID == "" {
 return nil
 }
 if err := db.Model(&models.Seat{}).Where("event_id = ?", eventID).Count(&totalSeats).Error; err != nil {
		return err
	}

	if err := db.Model(&models.Seat{}).Where("event_id = ? AND status = ?", eventID, models.StatusAvailable).Count(&availableSeats).Error; err != nil {
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
	case percentage > 0.1:
		newStatus = models.AvailabilityMedium
	default:
		newStatus = models.AvailabilityLow
	}

	return db.Model(&models.Event{}).Where("id = ?", eventID).Update("availability", newStatus).Error
}
