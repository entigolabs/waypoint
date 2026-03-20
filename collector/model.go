package collector

import (
	"fmt"

	"github.com/entigolabs/waypoint/client"
	"github.com/entigolabs/waypoint/internal/db"
)

func toCategory(c client.Category) db.Category {
	return db.Category{
		ID:          c.ID,
		Name:        c.Name,
		Description: c.Description,
		EmsIds:      c.EmsIds,
	}
}

func toEMSCategory(c client.EMSCategory) (db.EMSCategory, error) {
	id, err := db.ToUUID(c.ID)
	if err != nil {
		return db.EMSCategory{}, fmt.Errorf("cannot convert ID `%s` to UUID: %w", c.ID, err)
	}
	return db.EMSCategory{
		ID:   id,
		Name: c.Name,
	}, nil
}

func toEmsTheme(t client.EMSTheme) db.EmsTheme {
	return db.EmsTheme{
		Code:          t.Code,
		DatasetsCount: t.DatasetsCount,
	}
}

func toEmsThemeTranslation(t client.EMSThemeTranslation) db.EmsThemeTranslation {
	return db.EmsThemeTranslation{
		ID:          t.ID,
		Language:    t.Language,
		Value:       t.Value,
		Description: t.Description,
		CreatedAt:   db.ToTimestamptz(t.CreatedAt),
		UpdatedAt:   db.ToTimestamptz(t.UpdatedAt),
	}
}
