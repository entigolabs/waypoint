package server

import (
	"github.com/entigolabs/waypoint/internal/db"
	"github.com/google/uuid"
)

func toCategory(r db.Category) Category {
	return Category{
		Id:          r.ID,
		Name:        r.Name,
		Description: &r.Description,
		EmsIds:      r.EmsIds,
	}
}

func toEmsCategory(r db.EMSCategory) EmsCategory {
	return EmsCategory{
		Id:   uuid.UUID(r.ID.Bytes),
		Name: r.Name,
	}
}

func toEmsThemeTranslation(t db.EmsThemeTranslation) EmsThemeTranslation {
	return EmsThemeTranslation{
		Language:    t.Language,
		Value:       t.Value,
		Description: t.Description,
	}
}

func toEmsTheme(t db.EmsTheme, translations []EmsThemeTranslation, emsIDs []uuid.UUID) EmsTheme {
	return EmsTheme{
		Code:          t.Code,
		DatasetsCount: t.DatasetsCount,
		Translations:  translations,
		EmsIds:        emsIDs,
	}
}
