package server

import (
	"context"

	"github.com/entigolabs/waypoint/internal/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// dbReader is the read-only database interface used by service functions.
// *db.DB satisfies this interface.
type dbReader interface {
	GetCategories(ctx context.Context) ([]db.Category, error)
	GetEmsCategories(ctx context.Context) ([]db.EMSCategory, error)
	GetEmsThemes(ctx context.Context) ([]db.EmsTheme, error)
	GetEmsThemeTranslations(ctx context.Context) ([]db.EmsThemeTranslation, error)
	GetEmsThemeEmsCategories(ctx context.Context) ([]db.EmsThemeEmsCategoryRow, error)
}

func listCategories(ctx context.Context, database dbReader) ([]Category, error) {
	rows, err := database.GetCategories(ctx)
	if err != nil {
		return nil, err
	}
	data := make([]Category, len(rows))
	for i, r := range rows {
		data[i] = toCategory(r)
	}
	return data, nil
}

func listEmsCategories(ctx context.Context, database dbReader) ([]EmsCategory, error) {
	rows, err := database.GetEmsCategories(ctx)
	if err != nil {
		return nil, err
	}
	data := make([]EmsCategory, len(rows))
	for i, r := range rows {
		data[i] = toEmsCategory(r)
	}
	return data, nil
}

func listEmsThemes(ctx context.Context, database dbReader) ([]EmsTheme, error) {
	themes, err := database.GetEmsThemes(ctx)
	if err != nil {
		return nil, err
	}

	translations, err := database.GetEmsThemeTranslations(ctx)
	if err != nil {
		return nil, err
	}

	joins, err := database.GetEmsThemeEmsCategories(ctx)
	if err != nil {
		return nil, err
	}

	translationsByThemeID := make(map[pgtype.UUID][]EmsThemeTranslation, len(themes))
	for _, t := range translations {
		translationsByThemeID[t.EmsThemeId] = append(translationsByThemeID[t.EmsThemeId], toEmsThemeTranslation(t))
	}

	emsIDsByThemeID := make(map[pgtype.UUID][]uuid.UUID, len(themes))
	for _, j := range joins {
		emsIDsByThemeID[j.EmsThemeId] = append(emsIDsByThemeID[j.EmsThemeId], j.EmsCategoryID.Bytes)
	}

	data := make([]EmsTheme, len(themes))
	for i, t := range themes {
		transl := translationsByThemeID[t.ID]
		if transl == nil {
			transl = []EmsThemeTranslation{}
		}
		emsIDs := emsIDsByThemeID[t.ID]
		if emsIDs == nil {
			emsIDs = []uuid.UUID{}
		}
		data[i] = toEmsTheme(t, transl, emsIDs)
	}

	return data, nil
}
