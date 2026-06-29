package main

import (
	"brew-chatbot/config"
	"brew-chatbot/gemini"
	"brew-chatbot/handler"
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	// Keep main() minimal — just call run() and handle fatal errors
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	// 1. Load config
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}
	slog.Info("Config loaded successfully!")

	// 2. Create Gemini client
	geminiClient, err := gemini.NewClient(cfg.GeminiAPIKey)
	if err != nil {
		return fmt.Errorf("creating Gemini client: %w", err)
	}
	slog.Info("Gemini client ready!")

	// 3. Register routes
	mux := setupRoutes(geminiClient)

	// 4. Start the server
	addr := ":" + cfg.Port
	slog.Info("Server running on http://localhost" + addr)
	return http.ListenAndServe(addr, mux)
}

// setupRoutes wires all handlers to their routes and returns the mux
func setupRoutes(geminiClient *gemini.Client) *http.ServeMux {
	mux := http.NewServeMux()

	chatHandler := &handler.ChatHandler{Gemini: geminiClient}
	chatStreamHandler := &handler.ChatStreamHandler{Client: geminiClient}

	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/chat", chatHandler.Handle)
	mux.HandleFunc("POST /chat/stream", chatStreamHandler.ServeHTTP)

	return mux
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Brew Chatbot is alive!")
}