package server

import (
	"net/http"

	"github.com/lejeunel/go-image-annotator/adapters/api/json"
	p "github.com/lejeunel/go-image-annotator/adapters/api/json/group"
	"github.com/lejeunel/go-image-annotator/adapters/api/models"
	"github.com/lejeunel/go-image-annotator/use-cases/group/create"
	"github.com/lejeunel/go-image-annotator/use-cases/group/update"
)

func (s *Server) FindGroup(w http.ResponseWriter, r *http.Request, name string) {
	s.Group.Find.Execute(r.Context(), name, p.NewPresenter(w, s.Logger))
}

func (s *Server) CreateGroup(w http.ResponseWriter, r *http.Request) {
	body, ok := json.MustDecodeJSON[models.NewGroup](w, r)
	if !ok {
		return
	}

	s.Group.Create.Execute(r.Context(), create.Request{
		Name:        body.Name,
		Description: body.Description,
	}, p.NewPresenter(w, s.Logger))
}

func (s *Server) UpdateGroup(w http.ResponseWriter, r *http.Request, name string) {
	body, ok := json.MustDecodeJSON[models.UpdateGroup](w, r)
	if !ok {
		return
	}

	s.Group.Update.Execute(r.Context(), update.Request{
		Name: name, NewName: body.Name,
		NewDescription: body.Description,
	}, p.NewPresenter(w, s.Logger))
}

func (s *Server) DeleteGroup(w http.ResponseWriter, r *http.Request, name string) {
	s.Group.Delete.Execute(r.Context(), name, p.NewPresenter(w, s.Logger))
}

func (s *Server) ListGroups(w http.ResponseWriter, r *http.Request) {
	s.Group.List.Execute(r.Context(), p.NewPresenter(w, s.Logger))
}
