package db

import (
	"encoding/hex"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type Category struct {
	ID          int      `ksql:"id"`
	Name        string   `ksql:"name"`
	Description string   `ksql:"description"`
	EmsIds      []string `ksql:"ems_ids"`
}

type EMSCategory struct {
	ID   pgtype.UUID `ksql:"id"`
	Name string      `ksql:"name"`
}

type EmsTheme struct {
	ID            pgtype.UUID `ksql:"id"`
	Code          string      `ksql:"code"`
	DatasetsCount int         `ksql:"datasets_count"`
}

type EmsThemeTranslation struct {
	ID          int                `ksql:"id"`
	EmsThemeId  pgtype.UUID        `ksql:"ems_theme_id"`
	Language    string             `ksql:"language"`
	Value       string             `ksql:"value"`
	Description *string            `ksql:"description"`
	CreatedAt   pgtype.Timestamptz `ksql:"created_at"`
	UpdatedAt   pgtype.Timestamptz `ksql:"updated_at"`
}

type EmsThemeEmsCategoryRow struct {
	EmsThemeId    pgtype.UUID `ksql:"ems_theme_id"`
	EmsCategoryID pgtype.UUID `ksql:"ems_category_id"`
}

type CollectionMetadata struct {
	ID              int        `ksql:"id"`
	LastCollectedAt *time.Time `ksql:"last_collected_at"`
}

func ToTimestamptz(t *time.Time) pgtype.Timestamptz {
	if t == nil || t.IsZero() {
		return pgtype.Timestamptz{Valid: false}
	}
	return pgtype.Timestamptz{Time: *t, InfinityModifier: pgtype.Finite, Valid: true}
}

func ToUUID(src string) (pgtype.UUID, error) {
	buf, err := parseUUID(src)
	if err != nil {
		return pgtype.UUID{}, err
	}
	return pgtype.UUID{Bytes: buf, Valid: true}, nil
}

func parseUUID(src string) (dst [16]byte, err error) {
	switch len(src) {
	case 36:
		src = src[0:8] + src[9:13] + src[14:18] + src[19:23] + src[24:]
	case 32:
		// dashes already stripped, assume valid
	default:
		return dst, fmt.Errorf("cannot parse UUID %v", src)
	}

	buf, err := hex.DecodeString(src)
	if err != nil {
		return dst, err
	}

	copy(dst[:], buf)
	return dst, err
}
