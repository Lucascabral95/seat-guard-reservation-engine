package seeds

import (
	"fmt"
	"math/rand"
	"time"

	"booking-service/internal/models"
	"booking-service/internal/services"

	"gorm.io/gorm"
)

type seedEventConfig struct {
	Event       models.Event
	SoldPercent int
}

type seatSectionConfig struct {
	Section string
	Count   int
	Price   float64
	Prefix  string
}

func ResetAndSeed(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {

		if err := tx.Exec(`
            TRUNCATE TABLE seats, events, booking_orders, checkouts
            RESTART IDENTITY CASCADE;
        `).Error; err != nil {
			return err
		}

		baseDate := time.Date(2028, 1, 1, 0, 0, 0, 0, time.UTC)

		seedConfigs := []seedEventConfig{
			{
				Event: models.Event{
					Name: "Coldplay - Music of the Spheres II", Description: "El regreso triunfal a Buenos Aires.", Location: "Estadio River Plate",
					Date: baseDate.AddDate(0, 0, 10), Price: 1500000, Gender: "POP",
					PosterURL: "https://images.unsplash.com/photo-1493225255756-d9584f8606e9?q=80&w=800&auto=format&fit=crop",
				}, SoldPercent: 100,
			},
			{
				Event: models.Event{
					Name: "Taylor Swift - The Eras Tour Returns", Description: "La artista más grande del mundo vuelve a Argentina.", Location: "Estadio River Plate",
					Date: baseDate.AddDate(0, 1, 5), Price: 2000000, Gender: "POP",
					PosterURL: "https://images.unsplash.com/photo-1540039155733-5bb30b53aa14?q=80&w=800&auto=format&fit=crop",
				}, SoldPercent: 98,
			},
			{
				Event: models.Event{
					Name: "Dua Lipa - Radical Optimism", Description: "Presentando su nuevo álbum en un show único.", Location: "Campo Argentino de Polo",
					Date: baseDate.AddDate(0, 1, 20), Price: 1800000, Gender: "POP",
					PosterURL: "https://images.unsplash.com/photo-1533174072545-e8d4aa97edf9?q=80&w=800&auto=format&fit=crop",
				}, SoldPercent: 85,
			},
			{
				Event: models.Event{
					Name: "Bruno Mars Live", Description: "Funk, Soul y Pop en una noche mágica.", Location: "Estadio Único La Plata",
					Date: baseDate.AddDate(0, 2, 10), Price: 1900000, Gender: "POP",
					PosterURL: "https://images.unsplash.com/photo-1501386761578-eac5c94b800a?q=80&w=800&auto=format&fit=crop",
				}, SoldPercent: 60,
			},
			{
				Event: models.Event{
					Name: "Adele - One Night Only", Description: "La voz más potente llega por primera vez.", Location: "Estadio River Plate",
					Date: baseDate.AddDate(0, 2, 25), Price: 2500000, Gender: "POP",
					PosterURL: "https://images.unsplash.com/photo-1516280440614-6697288d5d38?q=80&w=800&auto=format&fit=crop",
				}, SoldPercent: 90,
			},
			{
				Event: models.Event{
					Name: "Harry Styles - Love On Tour", Description: "El ídolo británico regresa con su estilo único.", Location: "Estadio River Plate",
					Date: baseDate.AddDate(0, 3, 15), Price: 1600000, Gender: "POP",
					PosterURL: "https://images.unsplash.com/photo-1574158622682-e40e69881006?q=80&w=800&auto=format&fit=crop",
				}, SoldPercent: 40,
			},

			// --- ROCK ---
			{
				Event: models.Event{
					Name: "Metallica World Tour", Description: "Noche de metal puro.", Location: "Estadio Velez Sarsfield",
					Date: baseDate.AddDate(0, 4, 1), Price: 2000000, Gender: "ROCK",
					PosterURL: "https://images.unsplash.com/photo-1598387993441-a364f854c3e1?q=80&w=800&auto=format&fit=crop",
				}, SoldPercent: 50,
			},
			{
				Event: models.Event{
					Name: "Red Hot Chili Peppers", Description: "Funk Rock californiano al extremo.", Location: "Estadio River Plate",
					Date: baseDate.AddDate(0, 4, 15), Price: 1800000, Gender: "ROCK",
					PosterURL: "https://images.unsplash.com/photo-1459749411177-0473ef4884f3?q=80&w=800&auto=format&fit=crop",
				}, SoldPercent: 70,
			},
			{
				Event: models.Event{
					Name: "Foo Fighters", Description: "La leyenda del rock continúa.", Location: "Estadio Velez Sarsfield",
					Date: baseDate.AddDate(0, 5, 5), Price: 1700000, Gender: "ROCK",
					PosterURL: "https://images.unsplash.com/photo-1508973379184-7517410fb0bc?q=80&w=800&auto=format&fit=crop",
				}, SoldPercent: 30,
			},
			{
				Event: models.Event{
					Name: "La Renga", Description: "El banquete se sirve otra vez.", Location: "Estadio Único La Plata",
					Date: baseDate.AddDate(0, 5, 20), Price: 1200000, Gender: "ROCK",
					PosterURL: "https://images.unsplash.com/photo-1524368535928-5b5e00ddc76b?q=80&w=800&auto=format&fit=crop",
				}, SoldPercent: 95,
			},
			{
				Event: models.Event{
					Name: "AC/DC - Power Up", Description: "Alto voltaje en Buenos Aires.", Location: "Estadio River Plate",
					Date: baseDate.AddDate(0, 6, 10), Price: 2200000, Gender: "ROCK",
					PosterURL: "https://images.unsplash.com/photo-1470229722913-7ea05107f5c3?q=80&w=800&auto=format&fit=crop",
				}, SoldPercent: 10,
			},
			{
				Event: models.Event{
					Name: "Green Day", Description: "Punk Rock para saltar toda la noche.", Location: "Estadio Velez Sarsfield",
					Date: baseDate.AddDate(0, 6, 25), Price: 1600000, Gender: "ROCK",
					PosterURL: "https://images.unsplash.com/photo-1481595357464-307136eedf66?q=80&w=800&auto=format&fit=crop",
				}, SoldPercent: 25,
			},

			// --- ELECTRONICA ---
			{
				Event: models.Event{
					Name: "Hernan Cattaneo - Sunsetstrip", Description: "Progressive House al atardecer.", Location: "Campo Argentino de Polo",
					Date: baseDate.AddDate(0, 7, 0), Price: 2500000, Gender: "ELECTRONICA",
					PosterURL: "https://images.unsplash.com/photo-1571266028243-3716950639dd?q=80&w=800&auto=format&fit=crop",
				}, SoldPercent: 88,
			},
			{
				Event: models.Event{
					Name: "David Guetta", Description: "Los hits de la electrónica mundial.", Location: "Movistar Arena",
					Date: baseDate.AddDate(0, 7, 15), Price: 1800000, Gender: "ELECTRONICA",
					PosterURL: "https://images.unsplash.com/photo-1563841930606-67e2bce48b78?q=80&w=800&auto=format&fit=crop",
				}, SoldPercent: 55,
			},
			{
				Event: models.Event{
					Name: "Tiësto Live", Description: "La leyenda del trance y house.", Location: "Mandarine Park",
					Date: baseDate.AddDate(0, 8, 5), Price: 1900000, Gender: "ELECTRONICA",
					PosterURL: "https://images.unsplash.com/photo-1557787163-1635e2efb160?q=80&w=800&auto=format&fit=crop",
				}, SoldPercent: 30,
			},
			{
				Event: models.Event{
					Name: "Tomorrowland Presenta", Description: "Una noche mágica de EDM.", Location: "Parque de la Ciudad",
					Date: baseDate.AddDate(0, 8, 20), Price: 3000000, Gender: "ELECTRONICA",
					PosterURL: "https://images.unsplash.com/photo-1533174072545-e8d4aa97edf9?q=80&w=800&auto=format&fit=crop",
				}, SoldPercent: 99,
			},

			// --- JAZZ & TEATRO ---
			{
				Event: models.Event{
					Name: "Jazz en el Parque", Description: "Festival de Jazz al aire libre.", Location: "Bosques de Palermo",
					Date: baseDate.AddDate(0, 9, 1), Price: 500000, Gender: "JAZZ",
					PosterURL: "https://images.unsplash.com/photo-1511192336575-5a79af67a629?q=80&w=800&auto=format&fit=crop",
				}, SoldPercent: 15,
			},
			{
				Event: models.Event{
					Name: "Norah Jones", Description: "Una velada íntima de Jazz y Pop.", Location: "Teatro Gran Rex",
					Date: baseDate.AddDate(0, 9, 15), Price: 1200000, Gender: "JAZZ",
					PosterURL: "https://images.unsplash.com/photo-1514525253440-b393452e8d26?q=80&w=800&auto=format&fit=crop",
				}, SoldPercent: 75,
			},
			{
				Event: models.Event{
					Name: "El Fantasma de la Ópera", Description: "El clásico de Broadway.", Location: "Teatro Ópera",
					Date: baseDate.AddDate(0, 10, 5), Price: 1100000, Gender: "TEATRO",
					PosterURL: "https://images.unsplash.com/photo-1503095392237-fc550ccc92dc?q=80&w=800&auto=format&fit=crop",
				}, SoldPercent: 40,
			},
			{
				Event: models.Event{
					Name: "Fuerza Bruta Wayra", Description: "Teatro físico inmersivo.", Location: "Estadio Obras",
					Date: baseDate.AddDate(0, 10, 20), Price: 900000, Gender: "TEATRO",
					PosterURL: "https://images.unsplash.com/photo-1470225620780-dba8ba36b745?q=80&w=800&auto=format&fit=crop",
				}, SoldPercent: 60,
			},
			{
				Event: models.Event{
					Name: "Cirque du Soleil - OVO", Description: "La magia del circo llega a la ciudad.", Location: "Costanera Sur",
					Date: baseDate.AddDate(0, 11, 1), Price: 2800000, Gender: "TEATRO",
					PosterURL: "https://images.unsplash.com/photo-1535525153412-5a42439a210d?q=80&w=800&auto=format&fit=crop",
				}, SoldPercent: 80,
			},
			{
				Event: models.Event{
					Name: "Les Miserables", Description: "La revolución en el escenario.", Location: "Teatro Colón",
					Date: baseDate.AddDate(0, 11, 15), Price: 1500000, Gender: "TEATRO",
					PosterURL: "https://images.unsplash.com/photo-1460723237483-7a6dc9d0b212?q=80&w=800&auto=format&fit=crop",
				}, SoldPercent: 20,
			},

			// --- METAL & VARIOS ---
			{
				Event: models.Event{
					Name: "Iron Maiden", Description: "The Future Past Tour.", Location: "Estadio Huracán",
					Date: baseDate.AddDate(0, 12, 5), Price: 1900000, Gender: "METAL",
					PosterURL: "https://images.unsplash.com/photo-1621330383321-df621dc91459?q=80&w=800&auto=format&fit=crop",
				}, SoldPercent: 92,
			},
			{
				Event: models.Event{
					Name: "Megadeth", Description: "Thrash metal legendario.", Location: "Movistar Arena",
					Date:  baseDate.AddDate(1, 0, 10), // Enero 2029
					Price: 1600000, Gender: "METAL",
					PosterURL: "https://images.unsplash.com/photo-1574459891823-74b868625693?q=80&w=800&auto=format&fit=crop",
				}, SoldPercent: 35,
			},
			{
				Event: models.Event{
					Name: "Lollapalooza 2029", Description: "Tres días de pura música.", Location: "Hipódromo de San Isidro",
					Date:  baseDate.AddDate(1, 2, 15), // Marzo 2029
					Price: 5000000, Gender: "VARIOS",
					PosterURL: "https://images.unsplash.com/photo-1459749411177-0473ef4884f3?q=80&w=800&auto=format&fit=crop",
				}, SoldPercent: 70,
			},
			{
				Event: models.Event{
					Name: "Primavera Sound", Description: "El festival de Barcelona en Buenos Aires.", Location: "Parque de los Niños",
					Date:  baseDate.AddDate(1, 10, 10), // Nov 2029
					Price: 4200000, Gender: "VARIOS",
					PosterURL: "https://images.unsplash.com/photo-1506157786151-b8491531f063?q=80&w=800&auto=format&fit=crop",
				}, SoldPercent: 10,
			},
			{
				Event: models.Event{
					Name: "Disney On Ice", Description: "Magia sobre hielo para toda la familia.", Location: "Luna Park",
					Date: baseDate.AddDate(0, 6, 15), Price: 800000, Gender: "VARIOS",
					PosterURL: "https://images.unsplash.com/photo-1546707012-c46675f12716?q=80&w=800&auto=format&fit=crop",
				}, SoldPercent: 45,
			},
			{
				Event: models.Event{
					Name: "Hans Zimmer Live", Description: "Las mejores bandas sonoras de cine.", Location: "Movistar Arena",
					Date: baseDate.AddDate(1, 3, 20), Price: 2100000, Gender: "VARIOS",
					PosterURL: "https://images.unsplash.com/photo-1465847899078-b413929f7120?q=80&w=800&auto=format&fit=crop",
				}, SoldPercent: 95,
			},
			{
				Event: models.Event{
					Name: "Quilmes Rock", Description: "El festival de rock nacional más grande.", Location: "Tecnópolis",
					Date: baseDate.AddDate(1, 4, 10), Price: 1300000, Gender: "VARIOS",
					PosterURL: "https://images.unsplash.com/photo-1514525253440-b393452e8d26?q=80&w=800&auto=format&fit=crop",
				}, SoldPercent: 5,
			},
			{
				Event: models.Event{
					Name: "Eric Clapton", Description: "El dios de la guitarra.", Location: "Estadio Velez Sarsfield",
					Date: baseDate.AddDate(0, 8, 25), Price: 1800000, Gender: "ROCK",
					PosterURL: "https://images.unsplash.com/photo-1511671782779-c97d3d27a1d4?q=80&w=800&auto=format&fit=crop",
				}, SoldPercent: 65,
			},
		}

		seatCfg := []seatSectionConfig{
			{Section: "VIP", Count: 15, Price: 300, Prefix: "A"},
			{Section: "PLATEA", Count: 35, Price: 200, Prefix: "B"},
			{Section: "GENERAL", Count: 50, Price: 150, Prefix: "G"},
		}

		allSeats := make([]models.Seat, 0)
		r := rand.New(rand.NewSource(time.Now().UnixNano()))

		for _, cfg := range seedConfigs {
			if err := tx.Create(&cfg.Event).Error; err != nil {
				return err
			}

			eventSeats := buildSeatsForEvent(cfg.Event.ID, seatCfg, cfg.SoldPercent, r)
			allSeats = append(allSeats, eventSeats...)
		}

		if err := tx.CreateInBatches(&allSeats, 1000).Error; err != nil {
			return err
		}

		for _, cfg := range seedConfigs {
			services.UpdateEventAvailability(tx, cfg.Event.ID)
		}

		return nil
	})
}

func buildSeatsForEvent(eventID string, sections []seatSectionConfig, soldPercent int, r *rand.Rand) []models.Seat {
	out := make([]models.Seat, 0)

	for _, s := range sections {
		for i := 1; i <= s.Count; i++ {
			status := models.StatusAvailable
			if r.Intn(100) < soldPercent {
				status = models.StatusSold
			}

			out = append(out, models.Seat{
				Section:  s.Section,
				Number:   fmt.Sprintf("%s%d", s.Prefix, i),
				Price:    s.Price,
				Status:   status,
				EventID:  eventID,
				TicketID: nil,
				LockedBy: nil,
				LockedAt: nil,
			})
		}
	}
	return out
}
