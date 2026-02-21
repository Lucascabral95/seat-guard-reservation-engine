package database

import (
	"booking-service/internal/config"
	"booking-service/internal/models"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB(ctx context.Context, cfg *config.Config) *gorm.DB {
	dsn := cfg.DBUrl

	const maxRetries = 5
	retryDelay := 2 * time.Second

	var db *gorm.DB
	var err error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		log.Printf("   -> Intento %d/%d de conexión a la DB...", attempt, maxRetries)

		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			PrepareStmt: true,
			Logger:      getLogger(),
		})
		if err != nil {
			log.Printf("⚠️  gorm.Open falló: %v", err)
			waitOrExit(ctx, retryDelay)
			retryDelay *= 2
			continue
		}

		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalf("❌ Failed to get underlying sql.DB: %v", err)
		}

		sqlDB.SetMaxOpenConns(cfg.DbMaxOpenConns)
		sqlDB.SetMaxIdleConns(cfg.DbMaxIdleConns)
		sqlDB.SetConnMaxLifetime(cfg.DbConnMaxLifeTime)
		sqlDB.SetConnMaxIdleTime(cfg.DbConnMaxIdleTime)

		ctxPing, cancel := context.WithTimeout(ctx, 5*time.Second)
		pingErr := sqlDB.PingContext(ctxPing)
		cancel()

		if pingErr != nil {
			log.Printf("⚠️  Ping falló (intento %d/%d): %v", attempt, maxRetries, pingErr)
			waitOrExit(ctx, retryDelay)
			retryDelay *= 2
			continue
		}

		log.Println("✅ Conexión a la DB establecida")
		return db
	}

	log.Fatalf("❌ No se pudo conectar a la DB después de %d intentos: %v", maxRetries, err)
	return nil
}

// waitOrExit espera el delay o sale si el contexto se cancela
func waitOrExit(ctx context.Context, delay time.Duration) {
	select {
	case <-ctx.Done():
		log.Fatal("❌ Contexto cancelado durante el retry de DB")
	case <-time.After(delay):
	}
}

func CloseDB(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}
	return sqlDB.Close()
}

func getLogger() logger.Interface {
	return logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             2 * time.Second,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)
}

func RunMigrations(db *gorm.DB) error {
	log.Println("   -> Iniciando automigrate manual...")
	err := db.AutoMigrate(
		&models.Event{},
		&models.Seat{},
		&models.BookingOrder{},
		&models.Checkout{},
		&models.TicketPDF{},
	)
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}
	log.Println("✅ Migrations completed")
	return nil
}
