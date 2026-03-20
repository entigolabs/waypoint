package client

import "time"

type Response[T interface{}] struct {
	Data T `json:"data"`
}

type Category struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	EmsIds      []string `json:"emsIds"`
}

type EMSCategory struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type EMSTheme struct {
	Code          string                `json:"code"`
	Translations  []EMSThemeTranslation `json:"translations"`
	EmsIds        []string              `json:"emsIds"`
	DatasetsCount int                   `json:"datasetsCount"`
}

type EMSThemeTranslation struct {
	ID          int        `json:"id"`
	Language    string     `json:"language"`
	Value       string     `json:"value"`
	Description *string    `json:"description,omitempty"`
	CreatedAt   *time.Time `json:"createdAt,omitempty"`
	UpdatedAt   *time.Time `json:"updatedAt,omitempty"`
}
