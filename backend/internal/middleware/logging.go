package middleware

import (
	"encoding/json"
	"log/slog"
	"net"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/entigolabs/waypoint/internal/config"
	"github.com/go-chi/chi/v5"
)

type LoggingResponseWriter struct {
	http.ResponseWriter
	StatusCode   int
	ErrorType    string
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
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				err := json.NewEncoder(w).Encode(ErrorResponse{Errors: []EntigoError{internalError()}})
				if err != nil {
					slog.Error("failed to json encode panic response", "error", err)
				}
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
			duration := time.Since(start)
			attrs := []any{
				"http.request.method", r.Method,
				"url.path", r.URL.Path,
				"http.response.status_code", lrw.StatusCode,
				"client.address", clientIP(r),
				"http.request.duration_ms", duration.Milliseconds(),
				"url.scheme", scheme(r),
			}
			if route := getRoutePattern(r); route != "" {
				attrs = append(attrs, "http.route", route)
			}
			if ua := r.UserAgent(); ua != "" {
				attrs = append(attrs, "user_agent.original", ua)
			}
			if lrw.ErrorType != "" {
				attrs = append(attrs, "error.type", lrw.ErrorType)
			}
			if lrw.ErrorMessage != "" {
				attrs = append(attrs, "error.message", lrw.ErrorMessage)
			}
			logger(r.Context(), "", attrs...)
		}()
		next.ServeHTTP(lrw, r)
	})
}

func scheme(r *http.Request) string {
	if r.TLS != nil {
		return "https"
	}
	if s := r.Header.Get("X-Forwarded-Proto"); s != "" {
		return s
	}
	return "http"
}

func getRoutePattern(r *http.Request) string {
	if rctx := chi.RouteContext(r.Context()); rctx != nil {
		return rctx.RoutePattern()
	}
	return ""
}

func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		if ip, _, ok := strings.Cut(xff, ","); ok {
			return strings.TrimSpace(ip)
		}
		return xff
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
