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
	// Load .env file if it exists — for local development
	// In Docker/production, env vars are injected directly so this is safely ignored
	godotenv.Load()

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return Config{}, fmt.Errorf("GEMINI_API_KEY is not set in .env")
	}

	return Config{
		GeminiAPIKey: apiKey,
		Port: "8080",
	}, nil
}