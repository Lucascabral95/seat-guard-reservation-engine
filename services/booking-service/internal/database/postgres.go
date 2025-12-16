package database

import (
	"booking-service/internal/config"
	"booking-service/internal/models"

	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(cfg *config.Config) *gorm.DB {
	dsn := cfg.DBUrl
	if dsn == "" {
		if cfg.DBHost == "" || cfg.DBUser == "" || cfg.DBPassword == "" || cfg.DBName == "" {
			log.Fatal("❌ Database config missing: set DB_URL or set DB_HOST, DB_USER, DB_PASSWORD, DB_NAME (and optionally DB_PORT)")
		}
		dsn = fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=require",
			cfg.DBHost,
			cfg.DBUser,
			cfg.DBPassword,
			cfg.DBName,
			cfg.DBPort,
		)
	}
	log.Println("   -> Intentando abrir conexión GORM...")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt: false,
	})
	if err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}

	log.Println("   -> Conexión abierta, iniciando automigrate...")

	err = db.AutoMigrate(
		&models.Event{},
		&models.Seat{},
		&models.BookingOrder{},
	)
	if err != nil {
		log.Fatal("❌ Failed to migrate database:", err)
	}
	log.Println("✅ Migrations completed")

	// AGREGADO: VACIAR BOOKING_ORDERS
	if err := db.Exec("TRUNCATE TABLE booking_orders CASCADE").Error; err != nil {
		log.Printf("⚠️ Warning: Could not truncate booking_orders: %v", err)
	} else {
		log.Println("🧹 BookingOrder table cleared")
	}

	return db
}
