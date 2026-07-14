package handler

import (
	"brew-chatbot/internal/db"
	"brew-chatbot/internal/httputil"
	"brew-chatbot/internal/middleware"
	"net/http"

	"github.com/jackc/pgx/v5/pgtype"
)

type SessionHandler struct {
	Queries *db.Queries
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