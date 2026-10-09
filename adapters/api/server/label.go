package server

import (
	"net/http"

	"github.com/lejeunel/go-image-annotator/adapters/api/json"
	"github.com/lejeunel/go-image-annotator/adapters/api/models"
	pa "github.com/lejeunel/go-image-annotator/shared/pagination"
	"github.com/lejeunel/go-image-annotator/use-cases/label/create"
)

func (s *Server) FindLabel(w http.ResponseWriter, r *http.Request, name string) {
	s.Label.Find.Execute(r.Context(), name, json.NewLabelPresenter(w, s.Logger))
}

func (s *Server) CreateLabel(w http.ResponseWriter, r *http.Request) {
	body, ok := json.MustDecodeJSON[models.NewLabel](w, r)
	if !ok {
		return
	}

	req := create.Request{
		Name:        body.Name,
		Description: body.Description,
	}
	s.Label.Create.Execute(r.Context(), req, json.NewLabelPresenter(w, s.Logger))
}

func (s *Server) DeleteLabel(w http.ResponseWriter, r *http.Request, name string) {
	s.Label.Delete.Execute(r.Context(), name, json.NewLabelPresenter(w, s.Logger))
}

func (s *Server) ListLabels(w http.ResponseWriter, r *http.Request, params ListLabelsParams) {
	req := pa.NewPaginationParamsFromOptional(params.PageSize, params.Page, s.Label.DefaultPageSize)
	s.Label.List.Execute(r.Context(), req, json.NewLabelPresenter(w, s.Logger))
}
