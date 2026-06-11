package server

import (
	"net/http"

	"github.com/lejeunel/go-image-annotator/adapters/api/json"
	p "github.com/lejeunel/go-image-annotator/adapters/api/json/label"
	"github.com/lejeunel/go-image-annotator/adapters/api/models"
	"github.com/lejeunel/go-image-annotator/use-cases/label/create"
	"github.com/lejeunel/go-image-annotator/use-cases/label/delete"
	"github.com/lejeunel/go-image-annotator/use-cases/label/list"
	"github.com/lejeunel/go-image-annotator/use-cases/label/read"
)

func (s *Server) FindLabelByName(w http.ResponseWriter, r *http.Request, name string) {
	s.Label.Find.Execute(r.Context(),
		read.Request{Name: name},
		p.NewFindPresenter(w, s.Logger))
}
func (s *Server) CreateLabel(w http.ResponseWriter, r *http.Request) {
	body, ok := json.MustDecodeJSON[models.NewLabel](w, r)
	if !ok {
		return
	}
	req := create.Request{Name: body.Name}
	if body.Description != nil {
		req.Description = *body.Description
	}
	s.Label.Create.Execute(r.Context(), req, &p.Create{Writer: w})
}
func (s *Server) DeleteLabelByName(w http.ResponseWriter, r *http.Request, name string) {
	s.Label.Delete.Execute(r.Context(),
		delete.Request{Name: name},
		p.NewDeletePresenter(w, s.Logger))
}
func (s *Server) ListLabels(w http.ResponseWriter, r *http.Request, params ListLabelsParams) {
	req := list.Request{Page: 1, PageSize: s.Label.DefaultPageSize}
	if p := params.Page; p != nil {
		req.Page = *p
	}
	if p := params.PageSize; p != nil {
		req.PageSize = *p
	}
	s.Label.List.Execute(r.Context(), req, p.NewListPresenter(w, s.Logger))
}
