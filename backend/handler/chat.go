package handler

import (
	"brew-chatbot/gemini"
	"encoding/json"
	"net/http"
)

type ChatHandler struct {
	Gemini *gemini.Client
}

type chatRequest struct {
	History []gemini.Message `json:"history"`
	UserMessage string `json:"userMessage`
}

type chatResponse struct {
	Reply string `json:"reply"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func (h *ChatHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req chatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return 
	}

	if req.UserMessage == "" {
		writeError(w, "userMessage cannot be empty", http.StatusBadRequest)
		return 
	}

	reply, err := h.Gemini.Chat(r.Context(), req.History, req.UserMessage)
	if err != nil {
		writeError(w, "failed to get response from AI", http.StatusInternalServerError)
		return
	}

	writeJSON(w, chatResponse{Reply: reply}, http.StatusOK)
}

func writeJSON(w http.ResponseWriter, data any, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, message string, statusCode int) {
	writeJSON(w, errorResponse{Error: message}, statusCode)
}