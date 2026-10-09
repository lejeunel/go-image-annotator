package server

import (
	"net/http"

	"github.com/lejeunel/go-image-annotator/adapters/api/json"
	"github.com/lejeunel/go-image-annotator/adapters/api/models"
	"github.com/lejeunel/go-image-annotator/use-cases/role/create"
	"github.com/lejeunel/go-image-annotator/use-cases/role/update"
)

func (s *Server) FindRole(w http.ResponseWriter, r *http.Request, name string) {
	s.Role.Find.Execute(r.Context(), name, json.NewRolePresenter(w, s.Logger))
}

func (s *Server) CreateRole(w http.ResponseWriter, r *http.Request) {
	body, ok := json.MustDecodeJSON[models.NewRole](w, r)
	if !ok {
		return
	}

	req := create.Request{Name: body.Name, Description: body.Description}
	if body.Methods != nil {
		req.Methods = *body.Methods
	}

	s.Role.Create.Execute(r.Context(), req, json.NewRolePresenter(w, s.Logger))
}

func (s *Server) UpdateRole(w http.ResponseWriter, r *http.Request, name string) {
	body, ok := json.MustDecodeJSON[models.UpdateRole](w, r)
	if !ok {
		return
	}

	req := update.Request{
		Name:           name,
		NewName:        body.Name,
		NewDescription: body.Description,
	}
	if body.Methods != nil {
		req.NewMethods = *body.Methods
	}

	s.Role.Update.Execute(r.Context(), req, json.NewRolePresenter(w, s.Logger))
}

func (s *Server) DeleteRole(w http.ResponseWriter, r *http.Request, name string) {
	s.Role.Delete.Execute(r.Context(), name, json.NewRolePresenter(w, s.Logger))
}

func (s *Server) ListRoles(w http.ResponseWriter, r *http.Request) {
	s.Role.List.Execute(r.Context(), json.NewRolePresenter(w, s.Logger))
}
