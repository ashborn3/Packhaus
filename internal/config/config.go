package config

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadConfig() error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return nil
}
