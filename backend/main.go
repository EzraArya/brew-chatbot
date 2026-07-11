package main

import (
	"brew-chatbot/config"
	"brew-chatbot/gemini"
	"brew-chatbot/handler"
	"brew-chatbot/internal/middleware"
	"brew-chatbot/internal/db"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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

    pool, err := initDB(context.Background(), cfg.DatabaseURL)
    if err != nil {
        return err
    }
    defer pool.Close()

    geminiClient, err := initGemini(cfg.GeminiAPIKey)
    if err != nil {
        return err
    }

    queries := db.New(pool)
    mux := setupRoutes(geminiClient, queries)
    handler := middleware.Logging(middleware.BodyLimit(1<<20)(mux))
    return startServer(cfg.Port, handler)
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
func setupRoutes(geminiClient *gemini.Client, queries *db.Queries) *http.ServeMux {
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

func initDB(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
    // 1. Create the connection pool
    pool, err := pgxpool.New(ctx, databaseURL)
    if err != nil {
        return nil, fmt.Errorf("connecting to database: %w", err)
    }
    slog.Info("Database connected")

    // 2. Run migrations
    m, err := migrate.New("file://db/migrations", databaseURL)
    if err != nil {
        pool.Close() // clean up the pool before returning
        return nil, fmt.Errorf("creating migrator: %w", err)
    }
    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        pool.Close()
        return nil, fmt.Errorf("running migrations: %w", err)
    }
    slog.Info("Migrations applied")

    return pool, nil
}

func initGemini(apiKey string) (*gemini.Client, error) {
    client, err := gemini.NewClient(apiKey)
    if err != nil {
        return nil, fmt.Errorf("creating Gemini client: %w", err)
    }
    slog.Info("Gemini client ready")
    return client, nil
}