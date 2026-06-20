package main

import (
	"brew-chatbot/config"
	"brew-chatbot/gemini"
	"brew-chatbot/handler"
	"fmt"
	"net/http"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Println("failed to load config:", err)
		return
	}
	fmt.Println("Config loaded! API key starts with:", cfg.GeminiAPIKey[:8]+"...")

	geminiClient, err := gemini.NewClient(cfg.GeminiAPIKey)
	if err != nil {
		fmt.Println("Failed to create Gemini Client:", err)
		return 
	}
	fmt.Println("Gemini Client Ready")

	chatHandler := &handler.ChatHandler{Gemini: geminiClient}

	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/chat", chatHandler.Handle)
	
	fmt.Println("Server running on http://localhost:" + cfg.Port)
	http.ListenAndServe(":"+cfg.Port, nil)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Brew Chatbot is alive")
}