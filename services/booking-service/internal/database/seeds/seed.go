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
            TRUNCATE TABLE seats, events, booking_orders, checkouts, ticket_pdfs
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
					PosterURL: "https://res.cloudinary.com/dywcuco2r/image/upload/v1767480863/15-facts-about-coldplay-1689324991_xq7rft.jpg",
				}, SoldPercent: 100,
			},
			{
				Event: models.Event{
					Name: "Taylor Swift - The Eras Tour Returns", Description: "La artista más grande del mundo vuelve a Argentina.", Location: "Estadio River Plate",
					Date: baseDate.AddDate(0, 1, 5), Price: 2000000, Gender: "POP",
					PosterURL: "https://res.cloudinary.com/dywcuco2r/image/upload/v1767480948/Taylor-Swift-wallpaper-HD-photo-Eras-Tour-concert-1920-x-1080-pixels-laptop_johzfh.jpg",
				}, SoldPercent: 98,
			},
			{
				Event: models.Event{
					Name: "Dua Lipa - Radical Optimism", Description: "Presentando su nuevo álbum en un show único.", Location: "Campo Argentino de Polo",
					Date: baseDate.AddDate(0, 1, 20), Price: 1800000, Gender: "POP",
					PosterURL: "https://res.cloudinary.com/dywcuco2r/image/upload/v1767481761/dua-lipa-4k-we-re-good-teaser-2i3ukp7991k5c9bt_cb4w3k.jpg",
				}, SoldPercent: 85,
			},
			{
				Event: models.Event{
					Name: "Bruno Mars Live", Description: "Funk, Soul y Pop en una noche mágica.", Location: "Estadio Único La Plata",
					Date: baseDate.AddDate(0, 2, 10), Price: 1900000, Gender: "POP",
					PosterURL: "https://res.cloudinary.com/dywcuco2r/image/upload/v1767481468/bruno-mars-24k-versace-ikfjxojbdzo1mezc_prqkvg.jpg",
				}, SoldPercent: 60,
			},
			{
				Event: models.Event{
					Name: "Adele - One Night Only", Description: "La voz más potente llega por primera vez.", Location: "Estadio River Plate",
					Date: baseDate.AddDate(0, 2, 25), Price: 2500000, Gender: "POP",
					PosterURL: "https://res.cloudinary.com/dywcuco2r/image/upload/v1767481235/Adele-HD-Photos-03866_j6yoqd.jpg",
				}, SoldPercent: 90,
			},
			{
				Event: models.Event{
					Name: "Harry Styles - Love On Tour", Description: "El ídolo británico regresa con su estilo único.", Location: "Estadio River Plate",
					Date: baseDate.AddDate(0, 3, 15), Price: 1600000, Gender: "POP",
					PosterURL: "https://res.cloudinary.com/dywcuco2r/image/upload/v1767481277/PROD_Harry-styles_horizontal-banner_1643689181072_j40hp8.jpg",
				}, SoldPercent: 40,
			},

			// --- ROCK ---
			{
				Event: models.Event{
					Name: "Metallica World Tour", Description: "Noche de metal puro.", Location: "Estadio Velez Sarsfield",
					Date: baseDate.AddDate(0, 4, 1), Price: 2000000, Gender: "ROCK",
					PosterURL: "https://res.cloudinary.com/dywcuco2r/image/upload/v1767481204/metallica-1983-white-logo-897jxacb5usqbbmq_imm3mm.jpg",
				}, SoldPercent: 50,
			},
			{
				Event: models.Event{
					Name: "Red Hot Chili Peppers", Description: "Funk Rock californiano al extremo.", Location: "Estadio River Plate",
					Date: baseDate.AddDate(0, 4, 15), Price: 1800000, Gender: "ROCK",
					PosterURL: "https://res.cloudinary.com/dywcuco2r/image/upload/v1767481082/1282962-2560x1440-desktop-hd-red-hot-chili-peppers-background-image_xp5w5p.jpg",
				}, SoldPercent: 70,
			},
			{
				Event: models.Event{
					Name: "Foo Fighters", Description: "La leyenda del rock continúa.", Location: "Estadio Velez Sarsfield",
					Date: baseDate.AddDate(0, 5, 5), Price: 1700000, Gender: "ROCK",
					PosterURL: "https://res.cloudinary.com/dywcuco2r/image/upload/v1767481336/foo-fighters-wallpaper-preview_krhg6h.jpg",
				}, SoldPercent: 30,
			},
			{
				Event: models.Event{
					Name: "La Renga", Description: "El banquete se sirve otra vez.", Location: "Estadio Único La Plata",
					Date: baseDate.AddDate(0, 5, 20), Price: 1200000, Gender: "ROCK",
					PosterURL: "https://res.cloudinary.com/dywcuco2r/image/upload/v1767481494/20200502151435_la-renga-03_fb8xtp.jpg",
				}, SoldPercent: 95,
			},
			{
				Event: models.Event{
					Name: "AC/DC - Power Up", Description: "Alto voltaje en Buenos Aires.", Location: "Estadio River Plate",
					Date: baseDate.AddDate(0, 6, 10), Price: 2200000, Gender: "ROCK",
					PosterURL: "https://res.cloudinary.com/dywcuco2r/image/upload/v1767481589/maxresdefault_cjbhv7.jpg",
				}, SoldPercent: 10,
			},
			{
				Event: models.Event{
					Name: "Green Day", Description: "Punk Rock para saltar toda la noche.", Location: "Estadio Velez Sarsfield",
					Date: baseDate.AddDate(0, 6, 25), Price: 1600000, Gender: "ROCK",
					PosterURL: "https://res.cloudinary.com/dywcuco2r/image/upload/v1767480974/542030_xbyq2f.jpg",
				}, SoldPercent: 25,
			},

			// --- ELECTRONICA ---
			{
				Event: models.Event{
					Name: "Hernan Cattaneo - Sunsetstrip", Description: "Progressive House al atardecer.", Location: "Campo Argentino de Polo",
					Date: baseDate.AddDate(0, 7, 0), Price: 2500000, Gender: "ELECTRONICA",
					PosterURL: "https://res.cloudinary.com/dywcuco2r/image/upload/v1767481612/VD4PSPSTLVAGBO7BR5673LWZCU_gfsxwo.jpg",
				}, SoldPercent: 88,
			},
			{
				Event: models.Event{
					Name: "David Guetta", Description: "Los hits de la electrónica mundial.", Location: "Movistar Arena",
					Date: baseDate.AddDate(0, 7, 15), Price: 1800000, Gender: "ELECTRONICA",
					PosterURL: "https://res.cloudinary.com/dywcuco2r/image/upload/v1767481513/566901_j4yh2w.jpg",
				}, SoldPercent: 55,
			},
			{
				Event: models.Event{
					Name: "Tiësto Live", Description: "La leyenda del trance y house.", Location: "Mandarine Park",
					Date: baseDate.AddDate(0, 8, 5), Price: 1900000, Gender: "ELECTRONICA",
					PosterURL: "https://res.cloudinary.com/dywcuco2r/image/upload/v1767481536/be8142c76b72ec1417b52fe2a08c3e5b6a3353641e11e3a6c0aededebddf0e7a_ecctay.jpg",
				}, SoldPercent: 30,
			},
			{
				Event: models.Event{
					Name: "Tomorrowland Presenta", Description: "Una noche mágica de EDM.", Location: "Parque de la Ciudad",
					Date: baseDate.AddDate(0, 8, 20), Price: 3000000, Gender: "ELECTRONICA",
					PosterURL: "https://res.cloudinary.com/dywcuco2r/image/upload/v1767481004/tomorrowland-carnival-set-up-6zto43czxfims05e_vyzwjc.jpg",
				}, SoldPercent: 99,
			},

			// --- JAZZ & TEATRO ---
			{
				Event: models.Event{
					Name: "Jazz en el Parque", Description: "Festival de Jazz al aire libre.", Location: "Bosques de Palermo",
					Date: baseDate.AddDate(0, 9, 1), Price: 500000, Gender: "JAZZ",
					PosterURL: "https://res.cloudinary.com/dywcuco2r/image/upload/v1767481405/jazz-s6-1920x1080_lqd345.jpg",
				}, SoldPercent: 15,
			},
			{
				Event: models.Event{
					Name: "Norah Jones", Description: "Una velada íntima de Jazz y Pop.", Location: "Teatro Gran Rex",
					Date: baseDate.AddDate(0, 9, 15), Price: 1200000, Gender: "JAZZ",
					PosterURL: "https://res.cloudinary.com/dywcuco2r/image/upload/v1767481643/thumb-1920-9933_hp8cle.jpg",
				}, SoldPercent: 75,
			},
			{
				Event: models.Event{
					Name: "El Fantasma de la Ópera", Description: "El clásico de Broadway.", Location: "Teatro Ópera",
					Date: baseDate.AddDate(0, 10, 5), Price: 1100000, Gender: "TEATRO",
					PosterURL: "https://res.cloudinary.com/dywcuco2r/image/upload/v1767481677/9d758bd7f30a33a3c0d228a2198d571e5b50ddafdd3f62a3f5e8fe394a6e41df_xwddxo.jpg",
				}, SoldPercent: 40,
			},
			{
				Event: models.Event{
					Name: "Fuerza Bruta Wayra", Description: "Teatro físico inmersivo.", Location: "Estadio Obras",
					Date: baseDate.AddDate(0, 10, 20), Price: 900000, Gender: "TEATRO",
					PosterURL: "https://res.cloudinary.com/dywcuco2r/image/upload/v1767481420/maxresdefault_htxity.jpg",
				}, SoldPercent: 60,
			},
			{
				Event: models.Event{
					Name: "Cirque du Soleil - OVO", Description: "La magia del circo llega a la ciudad.", Location: "Costanera Sur",
					Date: baseDate.AddDate(0, 11, 1), Price: 2800000, Gender: "TEATRO",
					PosterURL: "https://res.cloudinary.com/dywcuco2r/image/upload/v1767481136/6ae628a435a9d8d4c1df27dcaa7b27dfe2c363aaea750dc12907e9cf8077c0cc._UR1920_1080__h1kynz.jpg",
				}, SoldPercent: 80,
			},
			{
				Event: models.Event{
					Name: "Les Miserables", Description: "La revolución en el escenario.", Location: "Teatro Colón",
					Date: baseDate.AddDate(0, 11, 15), Price: 1500000, Gender: "TEATRO",
					PosterURL: "https://res.cloudinary.com/dywcuco2r/image/upload/v1767481730/0d5719c549b70ac287b3ffc8322f036bb1a960fcfa77bb87b8ae7a7b7624459b_uoxpdl.jpg",
				}, SoldPercent: 20,
			},

			// --- METAL & VARIOS ---
			{
				Event: models.Event{
					Name: "Iron Maiden", Description: "The Future Past Tour.", Location: "Estadio Huracán",
					Date: baseDate.AddDate(0, 12, 5), Price: 1900000, Gender: "METAL",
					PosterURL: "https://res.cloudinary.com/dywcuco2r/image/upload/v1767481057/rtOMV2_zexk9t.jpg",
				}, SoldPercent: 92,
			},
			{
				Event: models.Event{
					Name: "Megadeth", Description: "Thrash metal legendario.", Location: "Movistar Arena",
					Date:  baseDate.AddDate(1, 0, 10), // Enero 2029
					Price: 1600000, Gender: "METAL",
					PosterURL: "https://res.cloudinary.com/dywcuco2r/image/upload/v1767481303/megadeth-sbzrdaipee20qjkw_mkomxt.jpg",
				}, SoldPercent: 35,
			},
			{
				Event: models.Event{
					Name: "Lollapalooza 2029", Description: "Tres días de pura música.", Location: "Hipódromo de San Isidro",
					Date:  baseDate.AddDate(1, 2, 15), // Marzo 2029
					Price: 5000000, Gender: "VARIOS",
					PosterURL: "https://res.cloudinary.com/dywcuco2r/image/upload/v1767481121/lollapalooza_pr0drq.webp",
				}, SoldPercent: 70,
			},
			{
				Event: models.Event{
					Name: "Primavera Sound", Description: "El festival de Barcelona en Buenos Aires.", Location: "Parque de los Niños",
					Date:  baseDate.AddDate(1, 10, 10), // Nov 2029
					Price: 4200000, Gender: "VARIOS",
					PosterURL: "https://res.cloudinary.com/dywcuco2r/image/upload/v1767481434/Portada-2022-Indie-Club-19801320-10_dk9zms.jpg",
				}, SoldPercent: 10,
			},
			{
				Event: models.Event{
					Name: "Disney On Ice", Description: "Magia sobre hielo para toda la familia.", Location: "Luna Park",
					Date: baseDate.AddDate(0, 6, 15), Price: 800000, Gender: "VARIOS",
					PosterURL: "https://res.cloudinary.com/dywcuco2r/image/upload/v1767481036/disney-on-ice-1391135_rtefly.jpg",
				}, SoldPercent: 45,
			},
			{
				Event: models.Event{
					Name: "Hans Zimmer Live", Description: "Las mejores bandas sonoras de cine.", Location: "Movistar Arena",
					Date: baseDate.AddDate(1, 3, 20), Price: 2100000, Gender: "VARIOS",
					PosterURL: "https://res.cloudinary.com/dywcuco2r/image/upload/v1767481662/308845533aa98ba01ec55277206c8d4d53f498fecb8b8d4ed7d48080c7a6b4a1._SX1080_FMjpg__ppgu7t.jpg",
				}, SoldPercent: 95,
			},
			{
				Event: models.Event{
					Name: "Quilmes Rock", Description: "El festival de rock nacional más grande.", Location: "Tecnópolis",
					Date: baseDate.AddDate(1, 4, 10), Price: 1300000, Gender: "VARIOS",
					PosterURL: "https://res.cloudinary.com/dywcuco2r/image/upload/v1767481317/quilmes-rock-2025jpg_y5uwxw.webp",
				}, SoldPercent: 5,
			},
			{
				Event: models.Event{
					Name: "Eric Clapton", Description: "El dios de la guitarra.", Location: "Estadio Velez Sarsfield",
					Date: baseDate.AddDate(0, 8, 25), Price: 1800000, Gender: "ROCK",
					PosterURL: "https://res.cloudinary.com/dywcuco2r/image/upload/v1767481702/VS67PF_lsmsvm.png",
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
