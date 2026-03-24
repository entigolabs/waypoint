package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/entigolabs/waypoint/internal/config"
)

type LoggingResponseWriter struct {
	http.ResponseWriter
	StatusCode   int
	ErrorMessage string
}

func (lrw *LoggingResponseWriter) WriteHeader(code int) {
	lrw.StatusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func RecovererMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				slog.ErrorContext(r.Context(), "request panicked",
					"panic", rec,
					"stack", string(debug.Stack()),
				)
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := slog.InfoContext
		if r.URL.Path == config.HealthCheckPath || r.URL.Path == config.MetricsPath {
			logger = slog.DebugContext
		}
		start := time.Now()
		lrw := &LoggingResponseWriter{ResponseWriter: w, StatusCode: http.StatusOK}
		defer func() {
			attrs := []any{
				"method", r.Method,
				"path", r.URL.Path,
				"status", lrw.StatusCode,
				"remote", r.RemoteAddr,
				"duration", time.Since(start),
			}
			if lrw.ErrorMessage != "" {
				attrs = append(attrs, "error", lrw.ErrorMessage)
			}
			logger(r.Context(), "", attrs...)
		}()
		next.ServeHTTP(lrw, r)
	})
}
