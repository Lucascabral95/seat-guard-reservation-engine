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

	AWSRegion   string
	SQSQueueUrl string

	StripeSecretKey string

	Smtp_Host string
	Smtp_Port string
	Smtp_User string
	Smtp_Pass string
	Smtp_From string
	Workers   string
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

		// Agregarles en vars de Terraform!!!! ⚠️
		Smtp_Host: getEnv("SMTP_HOST", "smtp.gmail.com"),
		Smtp_Port: getEnv("SMTP_PORT", "587"),
		Smtp_User: getEnv("SMTP_USER", ""),
		Smtp_Pass: getEnv("SMTP_PASS", ""),
		Smtp_From: getEnv("SMTP_FROM", getEnv("EMAIL_FROM", "")),
		Workers:   getEnv("WORKERS", "10"),
	}
}

func getEnv(key, defaultVal string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultVal
}
