package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port      string
	JWTSecret string

	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string

	DBUrl   string
	ENVMode string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	return &Config{
		Port:      getEnv("PORT", "8080"),
		JWTSecret: os.Getenv("JWT_SECRET"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBUser:     getEnv("DB_USER", "localhost"),
		DBPassword: getEnv("DB_PASSWORD", "5432"),
		DBName:     getEnv("DB_NAME", "5432"),
		DBPort:     getEnv("DB_PORT", "5432"),

		DBUrl:   getEnv("DB_URL", ""),
		ENVMode: getEnv("ENV_MODE", "development"),
	}
}

func getEnv(key, defaultVal string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultVal
}
