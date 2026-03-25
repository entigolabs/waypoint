package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/entigolabs/waypoint/client"
	"github.com/entigolabs/waypoint/collector"
	"github.com/entigolabs/waypoint/internal/config"
	"github.com/entigolabs/waypoint/internal/db"
	"github.com/entigolabs/waypoint/internal/middleware"
	"github.com/entigolabs/waypoint/internal/version"
	"github.com/entigolabs/waypoint/server"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"
	"github.com/vingarcia/ksql"
	kpgx "github.com/vingarcia/ksql/adapters/kpgx5"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal(err)
	}

	logger, closeLogger, err := config.NewLogger(cfg.LogConfig)
	if err != nil {
		log.Fatal(err)
	}
	slog.SetDefault(logger)

	runErr := run(cfg)
	if runErr != nil {
		slog.Error("fatal error", "error", runErr)
	}
	if closeLogger != nil {
		_ = closeLogger(context.Background())
	}
	if runErr != nil {
		os.Exit(1)
	}
}

func run(cfg config.Config) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	time.Local = time.UTC

	version.PrintVersion()
	slog.Debug("Debug enabled")

	terminated := make(chan os.Signal, 1)
	signal.Notify(terminated, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	rawDB, err := kpgx.New(ctx, getDBUri(cfg.DBConfig), ksql.Config{})
	if err != nil {
		return err
	}

	database := db.NewDB(rawDB)
	apiClient := client.NewApiClient(cfg.APIBaseURL, client.NewHttpClient(30*time.Second, 3, cfg.UserAgent))

	go collector.NewCollector(database, apiClient).Start(ctx)

	srv := &http.Server{
		Addr:    cfg.ServerAddr,
		Handler: newRouter(database, cfg),
	}

	go func() {
		slog.Info("server starting", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("failed to start server", "error", err)
			cancel()
		}
	}()

	select {
	case <-terminated:
	case <-ctx.Done():
	}
	return shutdown(srv, rawDB)
}

func getDBUri(cfg config.DBConfig) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)
}

func newRouter(database *db.DB, cfg config.Config) http.Handler {
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: cfg.AllowedOrigins,
		AllowedMethods: []string{http.MethodOptions, http.MethodGet},
	})

	publicRouter := chi.NewRouter()
	publicRouter.Use(corsHandler.Handler)
	handlerOptions := server.StrictHTTPServerOptions{
		RequestErrorHandlerFunc:  middleware.ErrorHandler,
		ResponseErrorHandlerFunc: middleware.ErrorHandler,
	}
	_ = server.HandlerFromMux(server.NewStrictHandlerWithOptions(server.NewServer(database), nil, handlerOptions), publicRouter)

	r := chi.NewRouter()
	r.Use(middleware.LoggingMiddleware)
	r.Use(middleware.RecovererMiddleware)
	r.Get(config.HealthCheckPath, server.NewHealthHandler(database))
	r.Handle(config.MetricsPath, promhttp.Handler())
	r.Mount("/api", publicRouter)
	return r
}

func shutdown(srv *http.Server, database ksql.DB) error {
	slog.Info("shutting down server")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("failed to shutdown server", "error", err)
	}
	if err := database.Close(); err != nil {
		return err
	}
	slog.Info("server shutdown")
	return nil
}
