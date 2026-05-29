package json

import (
	"encoding/json"
	"errors"
	"github.com/lejeunel/go-image-annotator/adapters/api/models"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/lejeunel/go-image-annotator/shared/pagination"
	"net/http"
)

func BuildPaginationResponse(p pagination.Pagination) models.Pagination {
	return models.Pagination{
		Page:       p.Page,
		PageSize:   p.PageSize,
		TotalItems: p.TotalRecords,
		TotalPages: p.TotalPages,
	}

}

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func WriteError(w http.ResponseWriter, status int, msg string) {
	WriteJSON(w, status, map[string]string{
		"error": msg,
	})
}

func DecodeJSON[T any](r *http.Request) (*T, error) {
	defer r.Body.Close()

	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return nil, err
	}
	return &v, nil
}

func MustDecodeJSON[T any](w http.ResponseWriter, r *http.Request) (*T, bool) {
	body, err := DecodeJSON[T](r)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request body")
		return nil, false
	}
	return body, true

}

func HTTPStatusCodeFromErr(err error) int {
	switch {
	case errors.Is(err, e.ErrDuplicate):
		return http.StatusConflict
	case errors.Is(err, e.ErrValidation):
		return http.StatusBadRequest
	case errors.Is(err, e.ErrDependency):
		return http.StatusFailedDependency
	case errors.Is(err, e.ErrNotFound):
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}

type ErrorPresenter struct {
	Writer http.ResponseWriter
}

func (p ErrorPresenter) Error(err error) {
	WriteError(p.Writer, HTTPStatusCodeFromErr(err), err.Error())
}
