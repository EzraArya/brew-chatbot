package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
	"log/slog"

	"brew-chatbot/gemini"
	"brew-chatbot/internal/httputil"
)

type ChatStreamHandler struct {
	Client *gemini.Client
}

func (h *ChatStreamHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httputil.WriteError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		httputil.WriteError(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	var req chatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteError(w, "invalid request body", http.StatusBadRequest)
		fmt.Fprintf(w, "data: [ERROR] BAD_REQUEST\n\n")
		return
	}

	if strings.TrimSpace(req.UserMessage) == "" {
		httputil.WriteError(w, "userMessage cannot be empty", http.StatusBadRequest)
		fmt.Fprintf(w, "data: [ERROR] VALIDATION_FAILED\n\n")
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancel()
	
	err := h.Client.ChatStream(ctx, req.History, req.UserMessage, func(chunk string) {
		cleanChunk := strings.ReplaceAll(chunk, "\n", "\\n")		
		fmt.Fprintf(w, "data: %s\n\n", cleanChunk)
		flusher.Flush()
	})

	if err != nil {
		fmt.Fprintf(w, "data: [ERROR] AI_SERVICE_FAILED\n\n")
		slog.Error("gemini stream failed", "error", err) 
		flusher.Flush()
		return
	}

	fmt.Fprintf(w, "data: [DONE]\n\n")
	flusher.Flush()
}