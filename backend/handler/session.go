package handler

import (
	"brew-chatbot/gemini"
	"brew-chatbot/internal/db"
	"brew-chatbot/internal/httputil"
	"brew-chatbot/internal/middleware"
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type SessionHandler struct {
	Queries *db.Queries
	Gemini *gemini.Client
}

func (h *SessionHandler) Create(w http.ResponseWriter, r *http.Request) {
	deviceID, ok := middleware.DeviceIDFromContext(r.Context())
	if !ok {
		httputil.WriteError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	session, err := h.Queries.CreateSession(r.Context(), deviceID)
	if err != nil {
		httputil.WriteError(w, "failed to create session", http.StatusInternalServerError)
		return
	}

	httputil.WriteJSON(w, session, http.StatusCreated)
}

func (h *SessionHandler) List(w http.ResponseWriter, r *http.Request) {
	deviceID, ok := middleware.DeviceIDFromContext(r.Context())
	if !ok {
		httputil.WriteError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	sessions, err := h.Queries.ListSessionsByDevice(r.Context(), deviceID)
	if err != nil {
		httputil.WriteError(w, "failed to list sessions", http.StatusInternalServerError)
		return
	}

	if sessions == nil {
		sessions = []db.ListSessionsByDeviceRow{}
	}

	httputil.WriteJSON(w, sessions, http.StatusOK)
}

func (h *SessionHandler) Get(w http.ResponseWriter, r *http.Request) {
	deviceID, ok := middleware.DeviceIDFromContext(r.Context())
	if !ok {
		httputil.WriteError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var sessionID pgtype.UUID
	if err := sessionID.Scan(r.PathValue("id")); err != nil {
		httputil.WriteError(w, "invalid session ID", http.StatusBadRequest)
		return
	}

	session, err := h.Queries.GetSession(r.Context(), db.GetSessionParams{
		ID: sessionID,
		DeviceID: deviceID,
	})
	if err != nil {
		httputil.WriteError(w, "session not found", http.StatusNotFound)
		return
	}

	httputil.WriteJSON(w, session, http.StatusOK)
}

func (h *SessionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	deviceID, ok := middleware.DeviceIDFromContext(r.Context())
	if !ok {
		httputil.WriteError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var sessionID pgtype.UUID
	if err := sessionID.Scan(r.PathValue("id")); err != nil {
		httputil.WriteError(w, "invalid session ID", http.StatusBadRequest)
		return
	}

	if err := h.Queries.DeleteSession(r.Context(), db.DeleteSessionParams{
		ID: sessionID,
		DeviceID: deviceID,
	}); err != nil {
		httputil.WriteError(w, "failed to delete session", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *SessionHandler) GenerateTitle(w http.ResponseWriter, r *http.Request) {
    deviceID, ok := middleware.DeviceIDFromContext(r.Context())
    if !ok {
        httputil.WriteError(w, "unauthorized", http.StatusUnauthorized)
        return
    }

    var sessionID pgtype.UUID
    if err := sessionID.Scan(r.PathValue("id")); err != nil {
        httputil.WriteError(w, "invalid session ID", http.StatusBadRequest)
        return
    }

    // Verify session belongs to this device
    if _, err := h.Queries.GetSession(r.Context(), db.GetSessionParams{
        ID:       sessionID,
        DeviceID: deviceID,
    }); err != nil {
        httputil.WriteError(w, "session not found", http.StatusNotFound)
        return
    }

    var body struct {
        Message string `json:"message"`
    }
    if err := json.NewDecoder(r.Body).Decode(&body); err != nil || strings.TrimSpace(body.Message) == "" {
        httputil.WriteError(w, "message is required", http.StatusBadRequest)
        return
    }

    ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
    defer cancel()

    title, err := h.Gemini.GenerateTitle(ctx, body.Message)
    if err != nil {
        httputil.WriteError(w, "failed to generate title", http.StatusInternalServerError)
        return
    }

    if _, err := h.Queries.UpdateSessionTitle(r.Context(), db.UpdateSessionTitleParams{
        ID:       sessionID,
        DeviceID: deviceID,
        Title:    title,
    }); err != nil {
        httputil.WriteError(w, "failed to save title", http.StatusInternalServerError)
        return
    }

    httputil.WriteJSON(w, map[string]string{"title": title}, http.StatusOK)
}