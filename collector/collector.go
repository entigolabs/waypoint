package collector

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/entigolabs/waypoint/client"
	"github.com/entigolabs/waypoint/internal/db"
)

// APIClient is the interface for fetching upstream data.
// *client.ApiClient satisfies this interface.
type APIClient interface {
	FetchCategories(ctx context.Context) ([]client.Category, error)
	FetchEmsCategories(ctx context.Context) ([]client.EMSCategory, error)
	FetchEmsThemes(ctx context.Context) ([]client.EMSTheme, error)
}

// Collector fetches data from the upstream API and stores it in the database.
type Collector struct {
	db     db.DBProvider
	client APIClient
}

// NewCollector creates a Collector using the provided DB service and API client.
func NewCollector(database db.DBProvider, apiClient APIClient) *Collector {
	return &Collector{db: database, client: apiClient}
}

// Start runs the collection loop: collects immediately then re-checks every hour.
func (c *Collector) Start(ctx context.Context) {
	c.runCollection(ctx)

	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			c.runCollection(ctx)
		}
	}
}

func (c *Collector) runCollection(ctx context.Context) {
	if err := c.tryCollect(ctx); err != nil {
		slog.Error("data collection failed", "error", err)
	}
}

// tryCollect acquires a row-level lock on collection_metadata and, if 24 hours
// have elapsed since the last successful collection, runs a full collect cycle.
func (c *Collector) tryCollect(ctx context.Context) error {
	return c.db.Transaction(ctx, func(txDB db.DBProvider) error {
		meta, err := txDB.GetCollectionMetadataForUpdate(ctx)
		if errors.Is(err, db.ErrLockNotAcquired) {
			// Another instance holds the lock — skip this cycle.
			slog.Debug("collection lock not acquired")
			return nil
		}
		if err != nil {
			return fmt.Errorf("failed to acquire collection lock: %w", err)
		}

		if meta.LastCollectedAt != nil && time.Since(*meta.LastCollectedAt) < 24*time.Hour {
			return nil
		}

		slog.Info("starting data collection")
		return c.collect(ctx, txDB)
	})
}

// collect fetches all three data sets and persists them atomically.
func (c *Collector) collect(ctx context.Context, database db.DBProvider) error {
	categories, err := c.client.FetchCategories(ctx)
	if err != nil {
		return err
	}

	emsCategories, err := c.client.FetchEmsCategories(ctx)
	if err != nil {
		return err
	}

	emsThemes, err := c.client.FetchEmsThemes(ctx)
	if err != nil {
		return err
	}

	if err = c.insertData(ctx, database, categories, emsCategories, emsThemes); err != nil {
		return err
	}

	slog.Info("data collection completed",
		"categories", len(categories),
		"ems_categories", len(emsCategories),
		"ems_themes", len(emsThemes),
	)
	return nil
}

// insertData removes previous data and inserts new data.
func (c *Collector) insertData(ctx context.Context, database db.DBProvider, categories []client.Category, emsCategories []client.EMSCategory, emsThemes []client.EMSTheme) error {
	if err := database.ClearData(ctx); err != nil {
		return err
	}

	for _, cat := range categories {
		if err := database.InsertCategory(ctx, new(toCategory(cat))); err != nil {
			return err
		}
	}

	for _, emsCat := range emsCategories {
		category, err := toEMSCategory(emsCat)
		if err != nil {
			slog.Error("failed to convert ems category", "ems_category", emsCat.Name, "error", err)
			continue
		}
		if err := database.InsertEmsCategory(ctx, &category); err != nil {
			return err
		}
	}

	for _, theme := range emsThemes {
		if err := c.insertEmsTheme(ctx, database, theme); err != nil {
			return err
		}
	}

	if err := database.UpdateCollectionMetadata(ctx); err != nil {
		return fmt.Errorf("failed to update collection metadata: %w", err)
	}

	return nil
}

func (c *Collector) insertEmsTheme(ctx context.Context, database db.DBProvider, theme client.EMSTheme) error {
	emsTheme := toEmsTheme(theme)
	if err := database.InsertEmsTheme(ctx, &emsTheme); err != nil {
		return err
	}

	for _, translation := range theme.Translations {
		emsTranslation := toEmsThemeTranslation(translation)
		emsTranslation.EmsThemeId = emsTheme.ID
		if err := database.InsertEmsThemeTranslation(ctx, &emsTranslation); err != nil {
			return err
		}
	}

	for _, emsID := range theme.EmsIds {
		emsId, err := db.ToUUID(emsID)
		if err != nil {
			slog.Error("failed to convert ems id to UUID", "ems_id", emsID, "error", err)
			continue
		}
		if err := database.InsertEmsThemeEmsCategory(ctx, emsTheme.ID, emsId); err != nil {
			return err
		}
	}
	return nil
}
