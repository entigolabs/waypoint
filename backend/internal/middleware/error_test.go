package middleware

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWriteErrorResponse_KnownEntigoError(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	ee := EntigoError{Status: http.StatusNotFound, Code: "NotFound", Message: "resource not found"}
	WriteErrorResponse(w, r, ee, "")

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var resp ErrorResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	require.Len(t, resp.Errors, 1)
	assert.Equal(t, "NotFound", resp.Errors[0].Code)
	assert.Equal(t, "resource not found", resp.Errors[0].Message)
}

func TestWriteErrorResponse_UnknownError(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	WriteErrorResponse(w, r, errors.New("something went wrong"), "inner detail")

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var resp ErrorResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	require.Len(t, resp.Errors, 1)
	assert.Equal(t, "InternalServerError", resp.Errors[0].Code)
	assert.Equal(t, http.StatusText(http.StatusInternalServerError), resp.Errors[0].Message)
}

func TestWriteErrorResponse_HeadMethod_KnownError(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodHead, "/", nil)

	ee := EntigoError{Status: http.StatusForbidden, Code: "Forbidden", Message: "access denied"}
	WriteErrorResponse(w, r, ee, "")

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestWriteErrorResponse_HeadMethod_UnknownError(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodHead, "/", nil)

	WriteErrorResponse(w, r, errors.New("db down"), "")

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestWriteErrorResponse_LoggingResponseWriter(t *testing.T) {
	rec := httptest.NewRecorder()
	lrw := &LoggingResponseWriter{ResponseWriter: rec}
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	WriteErrorResponse(lrw, r, errors.New("unexpected"), "my inner message")

	assert.Equal(t, "my inner message", lrw.ErrorMessage)
}

func TestWriteErrorResponse_LoggingResponseWriter_KnownError(t *testing.T) {
	rec := httptest.NewRecorder()
	lrw := &LoggingResponseWriter{ResponseWriter: rec}
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	ee := EntigoError{Status: http.StatusBadRequest, Code: "BadRequest", Message: "invalid input"}
	WriteErrorResponse(lrw, r, ee, "")

	assert.Equal(t, "invalid input", lrw.ErrorMessage)
}
