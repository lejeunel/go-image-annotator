package server

import (
	"github.com/lejeunel/go-image-annotator/adapters/api/json"
	p "github.com/lejeunel/go-image-annotator/adapters/api/json/user"
	"github.com/lejeunel/go-image-annotator/adapters/api/models"
	u "github.com/lejeunel/go-image-annotator/entities/user"
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
func (s *Server) WhoAmI(w http.ResponseWriter, r *http.Request) {
	user := u.IdentityFromContext(r.Context())

	if user != nil {
		json.WriteJSON(w, 200, UserIdentity{
			Id:      user.Id,
			Groups:  user.Groups,
			Roles:   user.Roles,
			IsAdmin: user.IsAdmin,
		})
		return
	}
	http.Error(w, "failed fetching user's identity", http.StatusBadRequest)
}
