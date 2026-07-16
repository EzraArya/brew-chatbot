package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"brew-chatbot/gemini"
	"brew-chatbot/internal/db"
	"brew-chatbot/internal/httputil"
	"brew-chatbot/internal/middleware"

	"github.com/jackc/pgx/v5/pgtype"
)

type ChatStreamHandler struct {
	Client *gemini.Client
	Queries *db.Queries
}

// chatRequest is the request body for the streaming chat endpoint
type chatRequest struct {
	UserMessage string `json:"userMessage"`
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

	var deviceID pgtype.UUID
	deviceID, ok = middleware.DeviceIDFromContext(r.Context())
	if !ok {
	    httputil.WriteError(w, "unauthorized", http.StatusUnauthorized)
	    return
	}
	
	var sessionID pgtype.UUID
	if err := sessionID.Scan(r.PathValue("id")); err != nil {
	    httputil.WriteError(w, "invalid session ID", http.StatusBadRequest)
	    return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancel()

	_, err := h.Queries.GetSession(ctx, db.GetSessionParams{
	    ID:       sessionID,
	    DeviceID: deviceID,
	})
	
	if err != nil {
	    fmt.Fprintf(w, "data: [ERROR] SESSION_NOT_FOUND\n\n")
	    flusher.Flush()
	    return
	}

	dbMessage, err := h.Queries.GetMessagesBySession(ctx, sessionID)
	if err != nil {
		fmt.Fprintf(w, "data: [ERROR] DB_FAILED\n\n")
		flusher.Flush()
		return
	}

	var history []gemini.Message
	for _, m := range dbMessage {
		history = append(history, gemini.Message{
			Role: m.Role,
			Content: m.Content,
		})
	}
	
	var fullReply strings.Builder

	err = h.Client.ChatStream(ctx, history, req.UserMessage, func(chunk string) {
	    fullReply.WriteString(chunk)  // accumulate
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

	if _, err := h.Queries.CreateMessage(ctx, db.CreateMessageParams{
	    SessionID: sessionID,
	    Role:      "user",
	    Content:   req.UserMessage,
	}); err != nil {
	    slog.Error("failed to save user message", "error", err)
	}
	
	if _, err := h.Queries.CreateMessage(ctx, db.CreateMessageParams{
	    SessionID: sessionID,
	    Role:      "model",
	    Content:   fullReply.String(),
	}); err != nil {
	    slog.Error("failed to save model message", "error", err)
	}
	
	flusher.Flush()
}