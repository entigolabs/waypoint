package server_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/entigolabs/waypoint/server"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// noopStrictServer implements StrictServerInterface with empty responses.
type noopStrictServer struct{}

func (noopStrictServer) GetCoreCategories(_ context.Context, _ server.GetCoreCategoriesRequestObject) (server.GetCoreCategoriesResponseObject, error) {
	return server.GetCoreCategories200JSONResponse(server.CategoriesResponse{
		Data:     []server.Category{},
		Metadata: server.Metadata{Total: 0},
	}), nil
}

func (noopStrictServer) GetCoreEmsCategories(_ context.Context, _ server.GetCoreEmsCategoriesRequestObject) (server.GetCoreEmsCategoriesResponseObject, error) {
	return server.GetCoreEmsCategories200JSONResponse(server.EmsCategoriesResponse{
		Data:     []server.EmsCategory{},
		Metadata: server.Metadata{Total: 0},
	}), nil
}

func (noopStrictServer) GetCoreEmsThemes(_ context.Context, _ server.GetCoreEmsThemesRequestObject) (server.GetCoreEmsThemesResponseObject, error) {
	return server.GetCoreEmsThemes200JSONResponse(server.EmsThemesResponse{
		Data:     []server.EmsTheme{},
		Metadata: server.Metadata{Total: 0},
	}), nil
}

func newTestRouter(ssi server.StrictServerInterface) *chi.Mux {
	opts := server.StrictHTTPServerOptions{
		RequestErrorHandlerFunc:  func(w http.ResponseWriter, r *http.Request, err error) { http.Error(w, err.Error(), 500) },
		ResponseErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) { http.Error(w, err.Error(), 500) },
	}
	r := chi.NewRouter()
	_ = server.HandlerFromMux(server.NewStrictHandlerWithOptions(ssi, nil, opts), r)
	return r
}

// TestCategoriesResponseShape verifies JSON marshalling of CategoriesResponse.
func TestCategoriesResponseShape(t *testing.T) {
	cr := server.CategoriesResponse{
		Data: []server.Category{
			{Id: 1, Name: "Science", EmsIds: []string{"1", "2"}},
		},
		Metadata: server.Metadata{Total: 1},
	}
	b, err := json.Marshal(cr)
	require.NoError(t, err)

	var m map[string]interface{}
	require.NoError(t, json.Unmarshal(b, &m))
	assert.Contains(t, m, "data")
	assert.Contains(t, m, "metadata")

	data := m["data"].([]interface{})
	require.Len(t, data, 1)
	item := data[0].(map[string]interface{})
	assert.Equal(t, float64(1), item["id"])
	assert.Equal(t, "Science", item["name"])
	emsIds := item["emsIds"].([]interface{})
	assert.Equal(t, []interface{}{"1", "2"}, emsIds)
}

// TestEmsCategoriesResponseShape verifies JSON marshalling of EmsCategoriesResponse.
func TestEmsCategoriesResponseShape(t *testing.T) {
	uid := uuid.MustParse("18fbac3f-2159-409e-a45c-a39ab9e085f8")
	r := server.EmsCategoriesResponse{
		Data: []server.EmsCategory{
			{Id: uid, Name: "ÜLDMÕISTED"},
		},
		Metadata: server.Metadata{Total: 1},
	}
	b, err := json.Marshal(r)
	require.NoError(t, err)

	var m map[string]interface{}
	require.NoError(t, json.Unmarshal(b, &m))

	data := m["data"].([]interface{})
	require.Len(t, data, 1)
	item := data[0].(map[string]interface{})
	assert.Equal(t, "ÜLDMÕISTED", item["name"])
	assert.Equal(t, "18fbac3f-2159-409e-a45c-a39ab9e085f8", item["id"])
}

// TestEmsThemesResponseShape verifies JSON marshalling of EmsThemesResponse.
func TestEmsThemesResponseShape(t *testing.T) {
	r := server.EmsThemesResponse{
		Data: []server.EmsTheme{
			{
				Code:          "TECH",
				DatasetsCount: 100,
				EmsIds:        []uuid.UUID{},
				Translations:  []server.EmsThemeTranslation{{Language: "en", Value: "Science"}},
			},
		},
		Metadata: server.Metadata{Total: 1},
	}
	b, err := json.Marshal(r)
	require.NoError(t, err)

	var m map[string]interface{}
	require.NoError(t, json.Unmarshal(b, &m))
	data := m["data"].([]interface{})
	item := data[0].(map[string]interface{})
	assert.Equal(t, "TECH", item["code"])
	assert.Equal(t, float64(100), item["datasetsCount"])
	transl := item["translations"].([]interface{})
	assert.Len(t, transl, 1)
	assert.Equal(t, "en", transl[0].(map[string]interface{})["language"])
}

// TestHandlerEndpoints verifies the router registers the three new endpoints and returns 200.
func TestHandlerEndpoints(t *testing.T) {
	r := newTestRouter(noopStrictServer{})

	for _, path := range []string{"/v1/core/categories", "/v1/core/ems-categories", "/v1/core/ems-themes"} {
		t.Run(path, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, path, nil)
			r.ServeHTTP(rec, req)
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

			var body map[string]interface{}
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
			assert.Contains(t, body, "data")
			assert.Contains(t, body, "metadata")
		})
	}
}
