package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"brew-chatbot/gemini"
	"brew-chatbot/internal/httputil"
)

type ChatStreamHandler struct {
	Client *gemini.Client
}

func (h *ChatStreamHandler) ServerHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httputil.WriteError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		httputil.WriteError(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	var req chatRequest
	if err := json.NewDecoder(r.Body).Decode(*&req); err != nil {
		httputil.WriteError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.UserMessage) == "" {
		httputil.WriteError(w, "userMessage cannot be empty", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	err := h.Client.ChatStream(r.Context(), req.History, req.UserMessage, func(chunk string) {
		cleanChunk := strings.ReplaceAll(chunk, "\n", "\\n")
		fmt.Fprintf(w, "data: %s\n\n", cleanChunk)
		flusher.Flush()
	})

	if err != nil {
		fmt.Fprintf(w, "data: [ERROR]\n\n")
		flusher.Flush()
		return
	}

	fmt.Fprintf(w, "data: [DONE]\n\n")
	flusher.Flush()
}