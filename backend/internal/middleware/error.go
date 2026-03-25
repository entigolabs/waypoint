package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
)

type ErrorResponse struct {
	Errors []EntigoError `json:"errors"`
}

type EntigoError struct {
	Status  int    `json:"-"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e EntigoError) Error() string {
	return fmt.Sprintf("code=%s, message=%s", e.Code, e.Message)
}

func ErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	WriteErrorResponse(w, r, err, "")
}

func internalError() EntigoError {
	return EntigoError{
		Status:  http.StatusInternalServerError,
		Code:    "InternalServerError",
		Message: http.StatusText(http.StatusInternalServerError),
	}
}

func WriteErrorResponse(w http.ResponseWriter, r *http.Request, err error, innerMessage string) {
	var ee EntigoError
	ok := errors.As(err, &ee)
	if r.Method == http.MethodHead {
		if ok {
			w.WriteHeader(ee.Status)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}

	if !ok {
		slog.Error(innerMessage, "error", err)
		ee = internalError()
	}

	if lrw, typeOk := w.(*LoggingResponseWriter); typeOk {
		if innerMessage != "" {
			lrw.ErrorMessage = innerMessage
		} else if ok {
			lrw.ErrorType = ee.Code
			lrw.ErrorMessage = ee.Message
		} else {
			lrw.ErrorMessage = err.Error()
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(ee.Status)
	err = json.NewEncoder(w).Encode(ErrorResponse{Errors: []EntigoError{ee}})
	if err != nil {
		slog.Error("failed to json encode error response", "error", err)
	}
}
