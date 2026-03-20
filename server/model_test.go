package server

import (
	"testing"

	"github.com/entigolabs/waypoint/internal/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
)

func TestToCategory(t *testing.T) {
	desc := "Test description"
	r := db.Category{ID: 1, Name: "Science", Description: desc, EmsIds: []string{"1", "2"}}
	cat := toCategory(r)
	assert.Equal(t, 1, cat.Id)
	assert.Equal(t, "Science", cat.Name)
	assert.Equal(t, &desc, cat.Description)
	assert.Equal(t, []string{"1", "2"}, cat.EmsIds)
}

func TestToEmsCategory(t *testing.T) {
	uid := [16]byte{1, 2, 3}
	r := db.EMSCategory{ID: pgtype.UUID{Bytes: uid, Valid: true}, Name: "ÜLDMÕISTED"}
	cat := toEmsCategory(r)
	assert.Equal(t, uuid.UUID(uid), cat.Id)
	assert.Equal(t, "ÜLDMÕISTED", cat.Name)
}

func TestToEmsThemeTranslation(t *testing.T) {
	desc := "A description"
	r := db.EmsThemeTranslation{Language: "en", Value: "Technology", Description: &desc}
	tr := toEmsThemeTranslation(r)
	assert.Equal(t, "en", tr.Language)
	assert.Equal(t, "Technology", tr.Value)
	assert.Equal(t, &desc, tr.Description)
}

func TestToEmsThemeTranslation_NilDescription(t *testing.T) {
	r := db.EmsThemeTranslation{Language: "et", Value: "Tehnoloogia", Description: nil}
	tr := toEmsThemeTranslation(r)
	assert.Nil(t, tr.Description)
}

func TestToEmsTheme(t *testing.T) {
	uid := [16]byte{5}
	theme := db.EmsTheme{ID: pgtype.UUID{Bytes: uid, Valid: true}, Code: "TECH", DatasetsCount: 99}
	translations := []EmsThemeTranslation{{Language: "en", Value: "Technology"}}
	emsIDs := []uuid.UUID{uuid.UUID(uid)}

	result := toEmsTheme(theme, translations, emsIDs)
	assert.Equal(t, "TECH", result.Code)
	assert.Equal(t, 99, result.DatasetsCount)
	assert.Equal(t, translations, result.Translations)
	assert.Equal(t, emsIDs, result.EmsIds)
}
