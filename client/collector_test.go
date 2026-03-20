package client_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/entigolabs/waypoint/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newApiClient(baseURL string) *client.ApiClient {
	return client.NewApiClient(baseURL, client.NewHttpClient(5*time.Second, 1))
}

func TestFetchCategories(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v2/core/categories", r.URL.Path)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{
				{"id": 1, "name": "Science", "description": "Science desc", "emsIds": []string{"1", "2"}},
				{"id": 2, "name": "Health", "description": "Health desc", "emsIds": []string{"3"}},
			},
		})
	}))
	defer srv.Close()

	items, err := newApiClient(srv.URL).FetchCategories(context.Background())
	require.NoError(t, err)
	require.Len(t, items, 2)
	assert.Equal(t, 1, items[0].ID)
	assert.Equal(t, "Science", items[0].Name)
	assert.Equal(t, []string{"1", "2"}, items[0].EmsIds)
}

func TestFetchEmsCategories(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v2/core/ems-categories", r.URL.Path)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{
				{"id": "18fbac3f-2159-409e-a45c-a39ab9e085f8", "name": "ÜLDMÕISTED"},
				{"id": "9d210571-5b9a-4778-b67c-91b7a9f82ca0", "name": "FILOSOOFIA"},
			},
		})
	}))
	defer srv.Close()

	items, err := newApiClient(srv.URL).FetchEmsCategories(context.Background())
	require.NoError(t, err)
	require.Len(t, items, 2)
	assert.Equal(t, "18fbac3f-2159-409e-a45c-a39ab9e085f8", items[0].ID)
	assert.Equal(t, "ÜLDMÕISTED", items[0].Name)
}

func TestFetchEmsThemes(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v2/core/ems-themes", r.URL.Path)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{
				{
					"code":          "TECH",
					"datasetsCount": 1059,
					"emsIds":        []string{"18fbac3f-2159-409e-a45c-a39ab9e085f8"},
					"translations": []map[string]any{
						{"id": 1, "language": "en", "value": "Science and Technology", "description": "Science and Technology", "createdAt": "2025-06-18T06:29:03.480Z", "updatedAt": "2025-06-18T06:29:03.480Z"},
						{"id": 2, "language": "et", "value": "Teadus ja tehnoloogia", "description": "Teadus ja tehnoloogia", "createdAt": "2025-06-18T06:29:03.480Z", "updatedAt": "2025-06-18T06:29:03.480Z"},
					},
				},
			},
		})
	}))
	defer srv.Close()

	items, err := newApiClient(srv.URL).FetchEmsThemes(context.Background())
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "TECH", items[0].Code)
	assert.Equal(t, 1059, items[0].DatasetsCount)
	assert.Len(t, items[0].Translations, 2)
	assert.Equal(t, "en", items[0].Translations[0].Language)
	assert.Len(t, items[0].EmsIds, 1)
}

func TestFetchCategoriesServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	_, err := newApiClient(srv.URL).FetchCategories(context.Background())
	require.Error(t, err)
}

func TestFetchEmsCategoriesServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer srv.Close()

	_, err := newApiClient(srv.URL).FetchEmsCategories(context.Background())
	require.Error(t, err)
}

func TestFetchEmsThemesServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
	}))
	defer srv.Close()

	_, err := newApiClient(srv.URL).FetchEmsThemes(context.Background())
	require.Error(t, err)
}
