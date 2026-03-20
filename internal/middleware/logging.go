package middleware

import (
	"log/slog"
	"net/http"
	"time"
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

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := &LoggingResponseWriter{ResponseWriter: w, StatusCode: http.StatusOK}
		defer func() {
			slog.Info("",
				"method", r.Method,
				"path", r.URL.Path,
				"status", lrw.StatusCode,
				"remote", r.RemoteAddr,
				"duration", time.Since(start),
				"error", lrw.ErrorMessage)
		}()
		next.ServeHTTP(lrw, r)
	})
}
