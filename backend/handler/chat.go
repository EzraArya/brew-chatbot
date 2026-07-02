package handler

import (
	"brew-chatbot/gemini"
	"brew-chatbot/internal/httputil"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

// ChatHandler holds the dependencies our handler needs
type ChatHandler struct {
	Gemini *gemini.Client
}

// chatRequest is what we expect iOS to send us
type chatRequest struct {
	History     []gemini.Message `json:"history"`
	UserMessage string           `json:"userMessage"`
}

// chatResponse is what we send back to iOS
type chatResponse struct {
	Reply string `json:"reply"`
}

// Handle is the HTTP handler for POST /chat
func (h *ChatHandler) Handle(w http.ResponseWriter, r *http.Request) {
	// 1. Only allow POST requests
	if r.Method != http.MethodPost {
		httputil.WriteError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 2. Decode the JSON body from iOS into our chatRequest struct
	var req chatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// 3. Basic validation
	if strings.TrimSpace(req.UserMessage) == "" {
		httputil.WriteError(w, "userMessage cannot be empty", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()
	
	// 4. Call Gemini
	reply, err := h.Gemini.Chat(ctx, req.History, req.UserMessage)
	if err != nil {
		httputil.WriteError(w, "failed to get response from AI", http.StatusInternalServerError)
		slog.Error("Gemini chat failed", "error", err)
		return
	}

	// 5. Send the response back to iOS as JSON
	httputil.WriteJSON(w, chatResponse{Reply: reply}, http.StatusOK)
}