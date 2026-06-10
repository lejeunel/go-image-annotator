package server

import (
	"github.com/lejeunel/go-image-annotator/adapters/api/json"
	p "github.com/lejeunel/go-image-annotator/adapters/api/json/user"
	"github.com/lejeunel/go-image-annotator/adapters/api/models"
	"github.com/lejeunel/go-image-annotator/use-cases/user/create"
	"net/http"
)

func (s *Server) CreateUser(w http.ResponseWriter, r *http.Request) {
	body, ok := json.MustDecodeJSON[models.NewUser](w, r)
	if !ok {
		return
	}

	req := create.Request{Id: body.Id}
	if *body.IsAdmin {
		req.IsAdmin = *body.IsAdmin
	}

	s.User.Create.Execute(
		r.Context(), req, p.NewCreatePresenter(w, s.Logger))

}
