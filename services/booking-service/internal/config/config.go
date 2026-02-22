package config

import (
	"log"
	"os"
	"strconv"
	"time"

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

	AWSRegion   string
	SQSQueueUrl string

	StripeSecretKey string

	Smtp_Host string
	Smtp_Port string
	Smtp_User string
	Smtp_Pass string
	Smtp_From string
	Workers   string

	DbMaxOpenConns    int
	DbMaxIdleConns    int
	DbConnMaxLifeTime time.Duration
	DbConnMaxIdleTime time.Duration
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

		AWSRegion:   getEnv("AWS_REGION", "us-east-1"),
		SQSQueueUrl: getEnv("SQS_QUEUE_URL", "sdsdsdsdsd"),

		StripeSecretKey: getEnv("STRIPE_SECRET_KEY", "sk_test_XXXXXXXXXXXXXXXXXXXX"),

		Smtp_Host: getEnv("SMTP_HOST", "smtp.gmail.com"),
		Smtp_Port: getEnv("SMTP_PORT", "587"),
		Smtp_User: getEnv("SMTP_USER", ""),
		Smtp_Pass: getEnv("SMTP_PASS", ""),
		Smtp_From: getEnv("SMTP_FROM", getEnv("EMAIL_FROM", "")),
		Workers:   getEnv("WORKERS", "10"),

		DbMaxOpenConns:    getEnvIntOrDefault("DB_MAX_OPEN_CONNS", 20),
		DbMaxIdleConns:    getEnvIntOrDefault("DB_MAX_IDLE_CONNS", 10),
		DbConnMaxLifeTime: getEnvDurationOrDefault("DB_CONN_MAX_LIFETIME", getEnvDurationOrDefault("DB_CONN_MAX_LIFE_TIME", 5*time.Minute)),
		DbConnMaxIdleTime: getEnvDurationOrDefault("DB_CONN_MAX_IDLE_TIME", 1*time.Minute),
	}
}

func getEnv(key, defaultVal string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultVal
}

func getEnvIntOrDefault(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	v, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Printf("Warning: invalid integer for %s (%s), using default %d", key, valueStr, defaultValue)
		return defaultValue
	}
	return v
}

func getEnvDurationOrDefault(key string, defaultValue time.Duration) time.Duration {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	v, err := time.ParseDuration(valueStr)
	if err != nil {
		log.Printf("Warning: invalid duration for %s (%s), using default %v", key, valueStr, defaultValue)
		return defaultValue
	}
	return v
}
