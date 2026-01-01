package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port     string
	IgUserID string
}

func LoadConfig() Config {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system env vars")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000" // default port
	}

	return Config{
		Port:     port,
		IgUserID: os.Getenv("IG_USER_ID"),
	}
}
