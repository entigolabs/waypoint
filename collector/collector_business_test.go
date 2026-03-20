package collector

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/entigolabs/waypoint/client"
	"github.com/entigolabs/waypoint/internal/db"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- mocks ---

type mockDB struct {
	lockErr         error
	meta            *db.CollectionMetadata
	clearCalled     bool
	metadataUpdated bool
}

func (m *mockDB) Transaction(_ context.Context, fn func(db.DBProvider) error) error {
	return fn(m)
}

func (m *mockDB) GetCollectionMetadataForUpdate(_ context.Context) (*db.CollectionMetadata, error) {
	return m.meta, m.lockErr
}

func (m *mockDB) UpdateCollectionMetadata(_ context.Context) error {
	m.metadataUpdated = true
	return nil
}

func (m *mockDB) ClearData(_ context.Context) error {
	m.clearCalled = true
	return nil
}

func (m *mockDB) InsertCategory(_ context.Context, _ *db.Category) error { return nil }

func (m *mockDB) InsertEmsCategory(_ context.Context, _ *db.EMSCategory) error { return nil }

func (m *mockDB) InsertEmsTheme(_ context.Context, _ *db.EmsTheme) error { return nil }

func (m *mockDB) InsertEmsThemeTranslation(_ context.Context, _ *db.EmsThemeTranslation) error {
	return nil
}

func (m *mockDB) InsertEmsThemeEmsCategory(_ context.Context, _, _ pgtype.UUID) error { return nil }

type mockAPIClient struct{}

func (m *mockAPIClient) FetchCategories(_ context.Context) ([]client.Category, error) {
	return []client.Category{}, nil
}

func (m *mockAPIClient) FetchEmsCategories(_ context.Context) ([]client.EMSCategory, error) {
	return []client.EMSCategory{}, nil
}

func (m *mockAPIClient) FetchEmsThemes(_ context.Context) ([]client.EMSTheme, error) {
	return []client.EMSTheme{}, nil
}

// newTestCollector wires a Collector with the given mock DB and a no-op API client.
func newTestCollector(database *mockDB) *Collector {
	return &Collector{db: database, client: &mockAPIClient{}}
}

// --- tests ---

func TestTryCollect_LockNotAcquired(t *testing.T) {
	database := &mockDB{lockErr: db.ErrLockNotAcquired}
	err := newTestCollector(database).tryCollect(context.Background())
	require.NoError(t, err)
	assert.False(t, database.clearCalled, "should not collect when lock is unavailable")
}

func TestTryCollect_RecentCollection(t *testing.T) {
	recent := time.Now().Add(-1 * time.Hour)
	database := &mockDB{meta: &db.CollectionMetadata{ID: 1, LastCollectedAt: &recent}}
	err := newTestCollector(database).tryCollect(context.Background())
	require.NoError(t, err)
	assert.False(t, database.clearCalled, "should not collect when last collection was recent")
	assert.False(t, database.metadataUpdated)
}

func TestTryCollect_NilLastCollectedAt(t *testing.T) {
	database := &mockDB{meta: &db.CollectionMetadata{ID: 1, LastCollectedAt: nil}}
	err := newTestCollector(database).tryCollect(context.Background())
	require.NoError(t, err)
	assert.True(t, database.clearCalled, "should collect when never collected before")
	assert.True(t, database.metadataUpdated)
}

func TestTryCollect_StaleCollection(t *testing.T) {
	old := time.Now().Add(-25 * time.Hour)
	database := &mockDB{meta: &db.CollectionMetadata{ID: 1, LastCollectedAt: &old}}
	err := newTestCollector(database).tryCollect(context.Background())
	require.NoError(t, err)
	assert.True(t, database.clearCalled, "should collect when last collection was over 24h ago")
	assert.True(t, database.metadataUpdated)
}

func TestTryCollect_LockError(t *testing.T) {
	database := &mockDB{lockErr: errors.New("connection lost")}
	err := newTestCollector(database).tryCollect(context.Background())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to acquire collection lock")
}
