package main

import (
	"brew-chatbot/config"
	"brew-chatbot/gemini"
	"brew-chatbot/handler"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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

    cfg, err := config.Load()
    if err != nil {
        return fmt.Errorf("loading config: %w", err)
    }
    slog.Info("Config loaded successfully!")

    geminiClient, err := gemini.NewClient(cfg.GeminiAPIKey)
    if err != nil {
        return fmt.Errorf("creating Gemini client: %w", err)
    }
    slog.Info("Gemini client ready!")

    mux := setupRoutes(geminiClient)
    return startServer(cfg.Port, mux)
}

func startServer(port string, handler http.Handler) error {
	server := &http.Server{
		Addr: ":" + port,
		Handler: handler,
		ReadTimeout: 15 * time.Second,
		WriteTimeout: 90 * time.Second,
		IdleTimeout: 120 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	go func() {
		slog.Info("Server running on http://localhost:" + port)
        if err := server.ListenAndServe(); err != http.ErrServerClosed {
            slog.Error("server error", "error", err)
        }
	}()

	<-ctx.Done()
	slog.Info("Shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("server shutdown: %w", err)
	}

	slog.Info("Server shut down cleanly")
	return nil
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