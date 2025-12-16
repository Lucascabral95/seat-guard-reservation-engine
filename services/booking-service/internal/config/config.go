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
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	return &Config{
		Port:      getEnv("PORT", "8080"),
		JWTSecret: os.Getenv("JWT_SECRET"),

		DBHost:     getEnv("DB_HOST", ""),
		DBUser:     getEnv("DB_USER", ""),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", ""),
		DBPort:     getEnv("DB_PORT", "5432"),

		DBUrl:   getEnv("DB_URL", ""),
		ENVMode: getEnv("ENV_MODE", "development"),

		AWSRegion:   getEnv("AWS_REGION", "us-east-1"),
		SQSQueueUrl: getEnv("SQS_QUEUE_URL", "sdsdsdsdsd"),

		StripeSecretKey: getEnv("STRIPE_SECRET_KEY", "sk_test_XXXXXXXXXXXXXXXXXXXX"),
	}
}

func getEnv(key, defaultVal string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultVal
}
