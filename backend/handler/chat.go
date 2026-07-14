package handler

import (
	"brew-chatbot/gemini"
	"brew-chatbot/internal/db"
	"brew-chatbot/internal/httputil"
	"brew-chatbot/internal/middleware"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// ChatHandler holds the dependencies our handler needs
type ChatHandler struct {
	Gemini  *gemini.Client
	Queries *db.Queries
}

// chatRequest is what we expect iOS to send us
type chatRequest struct {
	UserMessage string `json:"userMessage"`
}

// chatResponse is what we send back to iOS
type chatResponse struct {
	Reply string `json:"reply"`
}

// Handle is the HTTP handler for POST /sessions/{id}/chat
func (h *ChatHandler) Handle(w http.ResponseWriter, r *http.Request) {
	// 1. Only allow POST requests
	if r.Method != http.MethodPost {
		httputil.WriteError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 2. Decode the JSON body
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

	// 4. Get device ID from context (injected by DeviceID middleware)
	deviceID, ok := middleware.DeviceIDFromContext(r.Context())
	if !ok {
		httputil.WriteError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// 5. Parse session ID from URL path
	var sessionID pgtype.UUID
	if err := sessionID.Scan(r.PathValue("id")); err != nil {
		httputil.WriteError(w, "invalid session ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	// 6. Verify session belongs to this device
	if _, err := h.Queries.GetSession(ctx, db.GetSessionParams{
		ID:       sessionID,
		DeviceID: deviceID,
	}); err != nil {
		httputil.WriteError(w, "session not found", http.StatusNotFound)
		return
	}

	// 7. Fetch conversation history from DB
	dbMessages, err := h.Queries.GetMessagesBySession(ctx, sessionID)
	if err != nil {
		httputil.WriteError(w, "failed to load history", http.StatusInternalServerError)
		return
	}

	var history []gemini.Message
	for _, m := range dbMessages {
		history = append(history, gemini.Message{Role: m.Role, Content: m.Content})
	}

	// 8. Call Gemini
	reply, err := h.Gemini.Chat(ctx, history, req.UserMessage)
	if err != nil {
		httputil.WriteError(w, "failed to get response from AI", http.StatusInternalServerError)
		slog.Error("Gemini chat failed", "error", err)
		return
	}

	

	// Save user message
	if _, err := h.Queries.CreateMessage(ctx, db.CreateMessageParams{
	    SessionID: sessionID,
	    Role:      "user",
	    Content:   req.UserMessage,
	}); err != nil {
	    slog.Error("failed to save user message", "error", err)
	}
	
	// Save model reply
	if _, err := h.Queries.CreateMessage(ctx, db.CreateMessageParams{
	    SessionID: sessionID,
	    Role:      "model",
	    Content:   reply,
	}); err != nil {
	    slog.Error("failed to save model message", "error", err)
	}
	
	httputil.WriteJSON(w, chatResponse{Reply: reply}, http.StatusOK)
}