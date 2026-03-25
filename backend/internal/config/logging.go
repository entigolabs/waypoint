package config

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.40.0"
)

type LogLevel string
type LogOutput string
type LogFormat string

const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
)

const (
	LogOutputStdout LogOutput = "stdout"
	LogOutputFile   LogOutput = "file"
)

const (
	LogFormatText LogFormat = "text"
	LogFormatJSON LogFormat = "json"
	LogFormatOTEL LogFormat = "otel"
)

type levelHandler struct {
	level   slog.Leveler
	handler slog.Handler
}

func (h *levelHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level.Level()
}

func (h *levelHandler) Handle(ctx context.Context, r slog.Record) error {
	return h.handler.Handle(ctx, r)
}

func (h *levelHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &levelHandler{level: h.level, handler: h.handler.WithAttrs(attrs)}
}

func (h *levelHandler) WithGroup(name string) slog.Handler {
	return &levelHandler{level: h.level, handler: h.handler.WithGroup(name)}
}

func NewLogger(cfg LogConfig) (*slog.Logger, func(context.Context) error, error) {
	handler, shutdown, err := getHandler(cfg)
	if err != nil {
		return nil, nil, err
	}
	return slog.New(handler), shutdown, nil
}

func getOutputWriter(cfg LogConfig) io.Writer {
	switch cfg.LogOutput {
	case LogOutputFile:
		f, err := os.OpenFile(cfg.LogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err == nil {
			return f
		}
		slog.Warn("failed to open log file, defaulting to stdout", "error", err)
	case LogOutputStdout:
	default:
		slog.Warn("unknown log output, defaulting to stdout", "output", cfg.LogOutput)
	}
	return os.Stdout
}

func getHandler(cfg LogConfig) (slog.Handler, func(context.Context) error, error) {
	var level slog.Level
	_ = level.UnmarshalText([]byte(cfg.LogLevel))
	writer := getOutputWriter(cfg)

	opts := &slog.HandlerOptions{Level: level}
	if cfg.LogFormat == LogFormatText {
		return slog.NewTextHandler(writer, opts), nil, nil
	} else if cfg.LogFormat == LogFormatJSON {
		return slog.NewJSONHandler(writer, opts), nil, nil
	} else if cfg.LogFormat != LogFormatOTEL {
		slog.Warn("Log format is not supported, defaulting to otel")
	}
	return getOTELHandler(writer, level)
}

func getOTELHandler(writer io.Writer, level slog.Level) (slog.Handler, func(context.Context) error, error) {
	exporter, err := stdoutlog.New(stdoutlog.WithWriter(writer))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create log exporter: %w", err)
	}
	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("waypoint"),
		),
	)
	if err != nil {
		return nil, nil, err
	}
	lp := log.NewLoggerProvider(
		log.WithProcessor(log.NewSimpleProcessor(exporter)),
		log.WithResource(res),
	)
	handler := otelslog.NewHandler("waypoint-logger", otelslog.WithLoggerProvider(lp))
	return &levelHandler{level: level, handler: handler}, lp.Shutdown, nil
}
