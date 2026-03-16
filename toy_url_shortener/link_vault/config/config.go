package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strings"
)

type config struct {
	DbUrl        string
	RedisAddr    string
	KafkaBrokers []string
}

func Config() *config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	brokersStr := os.Getenv("KAFKA_BROKERS")
	if brokersStr == "" {
		log.Fatal("KAFKA_BROKERS environment variable not set")
	}

	// Split the string by comma to get a slice of strings
	kafkaBrokers := strings.Split(brokersStr, ",")

	config := config{DbUrl: os.Getenv("DB_URL"), RedisAddr: os.Getenv("REDIS_ADDR"), KafkaBrokers: kafkaBrokers}

	return &config
}
