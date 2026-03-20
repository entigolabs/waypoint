package server

import (
	"context"
	"errors"
	"testing"

	"github.com/entigolabs/waypoint/internal/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockDBReader struct {
	categories    []db.Category
	emsCategories []db.EMSCategory
	emsThemes     []db.EmsTheme
	translations  []db.EmsThemeTranslation
	joins         []db.EmsThemeEmsCategoryRow
	err           error
}

func (m *mockDBReader) GetCategories(_ context.Context) ([]db.Category, error) {
	return m.categories, m.err
}

func (m *mockDBReader) GetEmsCategories(_ context.Context) ([]db.EMSCategory, error) {
	return m.emsCategories, m.err
}

func (m *mockDBReader) GetEmsThemes(_ context.Context) ([]db.EmsTheme, error) {
	return m.emsThemes, m.err
}

func (m *mockDBReader) GetEmsThemeTranslations(_ context.Context) ([]db.EmsThemeTranslation, error) {
	return m.translations, m.err
}

func (m *mockDBReader) GetEmsThemeEmsCategories(_ context.Context) ([]db.EmsThemeEmsCategoryRow, error) {
	return m.joins, m.err
}

func TestListCategories(t *testing.T) {
	db := &mockDBReader{
		categories: []db.Category{
			{ID: 1, Name: "Science", Description: "desc", EmsIds: []string{"a"}},
			{ID: 2, Name: "Health", Description: "desc2", EmsIds: []string{}},
		},
	}
	result, err := listCategories(context.Background(), db)
	require.NoError(t, err)
	require.Len(t, result, 2)
	assert.Equal(t, 1, result[0].Id)
	assert.Equal(t, "Science", result[0].Name)
}

func TestListCategories_DBError(t *testing.T) {
	db := &mockDBReader{err: errors.New("connection lost")}
	_, err := listCategories(context.Background(), db)
	require.Error(t, err)
}

func TestListEmsCategories(t *testing.T) {
	uid := pgtype.UUID{Bytes: [16]byte{1}, Valid: true}
	db := &mockDBReader{
		emsCategories: []db.EMSCategory{
			{ID: uid, Name: "ÜLDMÕISTED"},
		},
	}
	result, err := listEmsCategories(context.Background(), db)
	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, uuid.UUID(uid.Bytes), result[0].Id)
	assert.Equal(t, "ÜLDMÕISTED", result[0].Name)
}

func TestListEmsThemes_GroupsByThemeID(t *testing.T) {
	themeID := pgtype.UUID{Bytes: [16]byte{1}, Valid: true}
	catID1 := pgtype.UUID{Bytes: [16]byte{2}, Valid: true}
	catID2 := pgtype.UUID{Bytes: [16]byte{3}, Valid: true}

	db := &mockDBReader{
		emsThemes: []db.EmsTheme{
			{ID: themeID, Code: "TECH", DatasetsCount: 5},
		},
		translations: []db.EmsThemeTranslation{
			{EmsThemeId: themeID, Language: "en", Value: "Technology"},
			{EmsThemeId: themeID, Language: "et", Value: "Tehnoloogia"},
		},
		joins: []db.EmsThemeEmsCategoryRow{
			{EmsThemeId: themeID, EmsCategoryID: catID1},
			{EmsThemeId: themeID, EmsCategoryID: catID2},
		},
	}

	result, err := listEmsThemes(context.Background(), db)
	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, "TECH", result[0].Code)
	assert.Equal(t, 5, result[0].DatasetsCount)
	assert.Len(t, result[0].Translations, 2)
	assert.Len(t, result[0].EmsIds, 2)
	assert.Contains(t, result[0].EmsIds, uuid.UUID(catID1.Bytes))
	assert.Contains(t, result[0].EmsIds, uuid.UUID(catID2.Bytes))
}

func TestListEmsThemes_EmptySlicesWhenNoRelatedData(t *testing.T) {
	themeID := pgtype.UUID{Bytes: [16]byte{1}, Valid: true}
	db := &mockDBReader{
		emsThemes: []db.EmsTheme{
			{ID: themeID, Code: "EMPTY"},
		},
	}

	result, err := listEmsThemes(context.Background(), db)
	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.NotNil(t, result[0].Translations)
	assert.NotNil(t, result[0].EmsIds)
	assert.Empty(t, result[0].Translations)
	assert.Empty(t, result[0].EmsIds)
}

func TestListEmsThemes_MultipleThemes(t *testing.T) {
	theme1 := pgtype.UUID{Bytes: [16]byte{1}, Valid: true}
	theme2 := pgtype.UUID{Bytes: [16]byte{2}, Valid: true}
	catID := pgtype.UUID{Bytes: [16]byte{9}, Valid: true}

	db := &mockDBReader{
		emsThemes: []db.EmsTheme{
			{ID: theme1, Code: "ALPHA"},
			{ID: theme2, Code: "BETA"},
		},
		translations: []db.EmsThemeTranslation{
			{EmsThemeId: theme2, Language: "en", Value: "Beta"},
		},
		joins: []db.EmsThemeEmsCategoryRow{
			{EmsThemeId: theme1, EmsCategoryID: catID},
		},
	}

	result, err := listEmsThemes(context.Background(), db)
	require.NoError(t, err)
	require.Len(t, result, 2)

	alpha := result[0]
	assert.Equal(t, "ALPHA", alpha.Code)
	assert.Empty(t, alpha.Translations)
	assert.Len(t, alpha.EmsIds, 1)

	beta := result[1]
	assert.Equal(t, "BETA", beta.Code)
	assert.Len(t, beta.Translations, 1)
	assert.Empty(t, beta.EmsIds)
}
