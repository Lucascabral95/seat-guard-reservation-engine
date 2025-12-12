package seeds

import (
	"fmt"
	"time"

	"booking-service/internal/models"

	"gorm.io/gorm"
)

type seatSectionConfig struct {
	Section string
	Count   int
	Price   float64
	Prefix  string
}

func ResetAndSeed(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// Reset solo estas tablas
		if err := tx.Exec(`
			TRUNCATE TABLE seats, events
			RESTART IDENTITY CASCADE;
		`).Error; err != nil {
			return err
		}

		now := time.Now()
		events := []models.Event{
			{Name: "Coldplay", Description: "Seed", Location: "Buenos Aires", Date: now.AddDate(0, 0, 7), Price: 15000},
			{Name: "Metallica", Description: "Seed", Location: "Buenos Aires", Date: now.AddDate(0, 0, 14), Price: 20000},
			{Name: "Dua Lipa", Description: "Seed", Location: "Buenos Aires", Date: now.AddDate(0, 0, 21), Price: 18000},
		}
		if err := tx.Create(&events).Error; err != nil {
			return err
		}

		cfg := []seatSectionConfig{
			{Section: "VIP", Count: 10, Price: 30000, Prefix: "A"},
			{Section: "PLATEA", Count: 20, Price: 20000, Prefix: "B"},
			{Section: "GENERAL", Count: 30, Price: 15000, Prefix: "G"},
		}

		allSeats := make([]models.Seat, 0, len(events)*(10+20+30))
		for _, e := range events {
			seats := buildSeatsForEvent(e.ID, cfg)
			allSeats = append(allSeats, seats...)
		}

		if err := tx.Create(&allSeats).Error; err != nil {
			return err
		}

		return nil
	})
}

func buildSeatsForEvent(eventID string, sections []seatSectionConfig) []models.Seat {
	out := make([]models.Seat, 0)
	for _, s := range sections {
		for i := 1; i <= s.Count; i++ {
			out = append(out, models.Seat{
				Section:  s.Section,
				Number:   fmt.Sprintf("%s%d", s.Prefix, i),
				Price:    s.Price,
				Status:   models.StatusAvailable,
				EventID:  eventID,
				TicketID: nil,
				LockedBy: nil,
				LockedAt: nil,
			})
		}
	}
	return out
}
