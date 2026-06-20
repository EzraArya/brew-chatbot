package main

import (
	"brew-chatbot/config"
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
	
	http.HandleFunc("/health", healthHandler)
	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Brew Chatbot is alive")
}