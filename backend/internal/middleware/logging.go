package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		reqID := uuid.New().String()

		ctx := context.WithValue(r.Context(), "request_id", reqID)
		r = r.WithContext(ctx)

		recorder := &statusRecorder{
			ResponseWriter: w,
			status: http.StatusOK,
		}

		next.ServeHTTP(recorder, r)
		duration := time.Since(start)

		slog.Info("HTTP Request",
					"method", r.Method,
					"path", r.URL.Path,
					"status", recorder.status,
					"duration", duration,
					"request_id", reqID,
				)
	})
}
