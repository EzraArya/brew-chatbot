package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	GeminiAPIKey string
	Port string
}

func Load() (Config, error) {
	err := godotenv.Load()
	if err != nil {
		return  Config{}, fmt.Errorf("error loading .env file: %w", err)
	}

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return Config{}, fmt.Errorf("GEMINI_API_KEY is not set in .env")
	}

	return Config{
		GeminiAPIKey: apiKey,
		Port: "8080",
	}, nil
}