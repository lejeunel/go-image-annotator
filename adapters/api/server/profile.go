package server

import (
	"net/http"

	"github.com/lejeunel/go-image-annotator/adapters/api/json"
	"github.com/lejeunel/go-image-annotator/adapters/api/models"
	"github.com/lejeunel/go-image-annotator/use-cases/profile/create"
	"github.com/lejeunel/go-image-annotator/use-cases/profile/update"
)

func (s *Server) FindProfile(w http.ResponseWriter, r *http.Request, name string) {
	s.Profile.Find.Execute(r.Context(), name, json.NewProfilePresenter(w, s.Logger))
}

func (s *Server) CreateProfile(w http.ResponseWriter, r *http.Request) {
	body, ok := json.MustDecodeJSON[models.NewProfile](w, r)
	if !ok {
		return
	}

	req := create.Request{
		Name:        body.Name,
		Description: body.Description,
		Group:       body.Group,
	}
	if body.Labels != nil {
		req.Labels = *body.Labels
	}

	s.Profile.Create.Execute(r.Context(), req, json.NewProfilePresenter(w, s.Logger))
}

func (s *Server) UpdateProfile(w http.ResponseWriter, r *http.Request, name string) {
	body, ok := json.MustDecodeJSON[models.UpdateProfile](w, r)
	if !ok {
		return
	}

	req := update.Request{
		Name:           name,
		NewName:        body.Name,
		NewDescription: body.Description,
		NewGroup:       body.Group,
	}
	if body.Labels != nil {
		req.NewLabels = *body.Labels
	}

	s.Profile.Update.Execute(r.Context(), req, json.NewProfilePresenter(w, s.Logger))
}

func (s *Server) DeleteProfile(w http.ResponseWriter, r *http.Request, name string) {
	s.Profile.Delete.Execute(r.Context(), name, json.NewProfilePresenter(w, s.Logger))
}
