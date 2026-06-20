package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/health", healthHandler)
	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Brew Chatbot is alive")
}