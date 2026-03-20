package server

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/entigolabs/waypoint/internal/db"
	"github.com/entigolabs/waypoint/internal/version"
)

type healthResponse struct {
	Name      string            `json:"name"`
	Version   string            `json:"version"`
	Status    string            `json:"status"`
	Time      time.Time         `json:"time"`
	StartTime time.Time         `json:"startTime"`
	BuildTime string            `json:"buildTime"`
	GitCommit string            `json:"gitCommit,omitempty"`
	Checks    map[string]string `json:"checks"`
}

// NewHealthHandler returns an HTTP handler for GET /health.
// It checks DB connectivity and reports name, version, build time, start time,
// and server time. Returns 200 when healthy, 503 when the database is unreachable.
func NewHealthHandler(database *db.DB) http.HandlerFunc {
	v := version.GetVersion()
	startTime := time.Now()
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()

		checks := map[string]string{}
		status := "ok"

		if err := database.Ping(ctx); err != nil {
			checks["database"] = "error"
			status = "error"
		} else {
			checks["database"] = "ok"
		}

		resp := healthResponse{
			Name:      "waypoint",
			Version:   v.Version,
			Status:    status,
			Time:      time.Now().UTC(),
			StartTime: startTime.UTC(),
			BuildTime: v.BuildDate,
			GitCommit: v.GitCommit,
			Checks:    checks,
		}

		statusCode := http.StatusOK
		if status != "ok" {
			statusCode = http.StatusServiceUnavailable
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		_ = json.NewEncoder(w).Encode(resp)
	}
}
