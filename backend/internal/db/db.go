package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/vingarcia/ksql"
)

var (
	categoryTable       = ksql.NewTable("core.categories", "id")
	emsCategoryTable    = ksql.NewTable("core.ems_categories", "id")
	emsThemesTable      = ksql.NewTable("core.ems_themes", "id")
	emsTranslationTable = ksql.NewTable("core.ems_theme_translations", "id")
)

// ErrLockNotAcquired is returned when another instance already holds the
// collection lock (SKIP LOCKED returned no rows).
var ErrLockNotAcquired = errors.New("collection lock not available")

// DBProvider is the interface the collector uses to interact with the database.
// *DB implements DBProvider.
type DBProvider interface {
	Transaction(ctx context.Context, fn func(DBProvider) error) error
	GetCollectionMetadataForUpdate(ctx context.Context) (*CollectionMetadata, error)
	UpdateCollectionMetadata(ctx context.Context) error
	ClearData(ctx context.Context) error
	InsertCategory(ctx context.Context, category *Category) error
	InsertEmsCategory(ctx context.Context, category *EMSCategory) error
	InsertEmsTheme(ctx context.Context, emsTheme *EmsTheme) error
	InsertEmsThemeTranslation(ctx context.Context, translation *EmsThemeTranslation) error
	InsertEmsThemeEmsCategory(ctx context.Context, themeId, emsCategoryID pgtype.UUID) error
}

// DB is the database service used by both the server and the collector.
type DB struct {
	db ksql.Provider
}

var _ DBProvider = (*DB)(nil)

// NewDB wraps a ksql.Provider in a DB service.
func NewDB(db ksql.Provider) *DB {
	return &DB{db: db}
}

// Transaction runs fn inside a database transaction. The DBProvider passed to
// fn is scoped to that transaction, so all operations within fn are atomic.
func (d *DB) Transaction(ctx context.Context, fn func(DBProvider) error) error {
	return d.db.Transaction(ctx, func(tx ksql.Provider) error {
		return fn(&DB{db: tx})
	})
}

// Ping executes a trivial query to verify the database connection is alive.
func (d *DB) Ping(ctx context.Context) error {
	_, err := d.db.Exec(ctx, "SELECT 1")
	return err
}

// GetCollectionMetadataForUpdate acquires an ACCESS EXCLUSIVE table-level lock on
// core.collection_metadata for the duration of the calling transaction, then reads
// the metadata row. NOWAIT means any competing instance receives an immediate
// PostgreSQL error (55P03) which is translated to ErrLockNotAcquired so the
// caller can skip cleanly. Only one instance can hold this lock at a time, so
// the collection is guaranteed to be single-threaded across all replicas.
func (d *DB) GetCollectionMetadataForUpdate(ctx context.Context) (*CollectionMetadata, error) {
	if _, err := d.db.Exec(ctx,
		"LOCK TABLE core.collection_metadata IN ACCESS EXCLUSIVE MODE NOWAIT",
	); err != nil {
		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok && pgErr.Code == "55P03" {
			return nil, ErrLockNotAcquired
		}
		return nil, err
	}

	var row CollectionMetadata
	if err := d.db.QueryOne(ctx, &row,
		"SELECT id, last_collected_at FROM core.collection_metadata WHERE id = 1",
	); err != nil {
		return nil, err
	}
	return &row, nil
}

// UpdateCollectionMetadata sets last_collected_at to the current time.
func (d *DB) UpdateCollectionMetadata(ctx context.Context) error {
	_, err := d.db.Exec(ctx,
		"UPDATE core.collection_metadata SET last_collected_at = NOW() WHERE id = 1")
	return err
}

// ClearData removes all collected data in FK-safe order.
func (d *DB) ClearData(ctx context.Context) error {
	for _, q := range []string{
		"DELETE FROM core.ems_theme_ems_categories",
		"DELETE FROM core.ems_theme_translations",
		"DELETE FROM core.ems_themes",
		"DELETE FROM core.ems_categories",
		"DELETE FROM core.categories",
	} {
		if _, err := d.db.Exec(ctx, q); err != nil {
			return fmt.Errorf("failed to clear data: %w", err)
		}
	}
	return nil
}

// InsertCategory inserts a single category row.
func (d *DB) InsertCategory(ctx context.Context, category *Category) error {
	if err := d.db.Insert(ctx, categoryTable, category); err != nil {
		return fmt.Errorf("failed to insert category %d: %w", category.ID, err)
	}
	return nil
}

// InsertEmsCategory inserts a single EMS category row.
func (d *DB) InsertEmsCategory(ctx context.Context, category *EMSCategory) error {
	if err := d.db.Insert(ctx, emsCategoryTable, category); err != nil {
		return fmt.Errorf("failed to insert ems_category %s: %w", category.ID, err)
	}
	return nil
}

// InsertEmsTheme inserts a single EMS theme row.
func (d *DB) InsertEmsTheme(ctx context.Context, theme *EmsTheme) error {
	if err := d.db.Insert(ctx, emsThemesTable, theme); err != nil {
		return fmt.Errorf("failed to insert ems_theme %s: %w", theme.Code, err)
	}
	return nil
}

// InsertEmsThemeTranslation inserts a single EMS theme translation row.
func (d *DB) InsertEmsThemeTranslation(ctx context.Context, translation *EmsThemeTranslation) error {
	if err := d.db.Insert(ctx, emsTranslationTable, translation); err != nil {
		return fmt.Errorf("failed to insert translation %d for theme %s: %w", translation.ID, translation.EmsThemeId, err)
	}
	return nil
}

// InsertEmsThemeEmsCategory inserts a join row linking an EMS theme to an EMS category.
func (d *DB) InsertEmsThemeEmsCategory(ctx context.Context, themeID, emsCategoryID pgtype.UUID) error {
	_, err := d.db.Exec(ctx,
		`INSERT INTO core.ems_theme_ems_categories (ems_theme_id, ems_category_id)
		 VALUES ($1, $2)`,
		themeID, emsCategoryID)
	if err != nil {
		return fmt.Errorf("failed to insert theme-category link %s→%s: %w", themeID, emsCategoryID, err)
	}
	return nil
}

// GetCategories returns all categories ordered by id.
func (d *DB) GetCategories(ctx context.Context) ([]Category, error) {
	var rows []Category
	if err := d.db.Query(ctx, &rows,
		"SELECT id, name, description, ems_ids FROM core.categories ORDER BY id",
	); err != nil {
		return nil, fmt.Errorf("failed to query categories: %w", err)
	}
	return rows, nil
}

// GetEmsCategories returns all EMS categories ordered by name.
func (d *DB) GetEmsCategories(ctx context.Context) ([]EMSCategory, error) {
	var rows []EMSCategory
	if err := d.db.Query(ctx, &rows,
		"SELECT id, name FROM core.ems_categories ORDER BY name",
	); err != nil {
		return nil, fmt.Errorf("failed to query ems_categories: %w", err)
	}
	return rows, nil
}

// GetEmsThemes returns all EMS themes ordered by code.
func (d *DB) GetEmsThemes(ctx context.Context) ([]EmsTheme, error) {
	var rows []EmsTheme
	if err := d.db.Query(ctx, &rows,
		"SELECT id, code, datasets_count FROM core.ems_themes ORDER BY code",
	); err != nil {
		return nil, fmt.Errorf("failed to query ems_themes: %w", err)
	}
	return rows, nil
}

// GetEmsThemeTranslations returns all EMS theme translations.
func (d *DB) GetEmsThemeTranslations(ctx context.Context) ([]EmsThemeTranslation, error) {
	var rows []EmsThemeTranslation
	if err := d.db.Query(ctx, &rows,
		"SELECT id, ems_theme_id, language, value, description FROM core.ems_theme_translations",
	); err != nil {
		return nil, fmt.Errorf("failed to query ems_theme_translations: %w", err)
	}
	return rows, nil
}

// GetEmsThemeEmsCategories returns all EMS theme–category join rows.
func (d *DB) GetEmsThemeEmsCategories(ctx context.Context) ([]EmsThemeEmsCategoryRow, error) {
	var rows []EmsThemeEmsCategoryRow
	if err := d.db.Query(ctx, &rows,
		"SELECT ems_theme_id, ems_category_id FROM core.ems_theme_ems_categories",
	); err != nil {
		return nil, fmt.Errorf("failed to query ems_theme_ems_categories: %w", err)
	}
	return rows, nil
}
