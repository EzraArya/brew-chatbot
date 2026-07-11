package middleware

import (
    "brew-chatbot/internal/httputil"
    "context"
    "net/http"

    "github.com/jackc/pgx/v5/pgtype"
)

type contextKey string
const deviceIDKey contextKey = "device_id"

func DeviceID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headerVal := r.Header.Get("X-Device-ID")

		if headerVal == "" {
		    httputil.WriteError(w, "X-Device-ID header is required", http.StatusUnauthorized)
		    return
		}

		var deviceID pgtype.UUID
		if err := deviceID.Scan(headerVal); err != nil {
			httputil.WriteError(w, "X-Device-ID must be a valid UUID", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), deviceIDKey, deviceID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func DeviceIDFromContext(ctx context.Context) (pgtype.UUID, bool) {
    id, ok := ctx.Value(deviceIDKey).(pgtype.UUID)
    return id, ok
}