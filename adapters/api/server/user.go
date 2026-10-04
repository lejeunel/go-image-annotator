package server

import (
	"net/http"

	"github.com/lejeunel/go-image-annotator/adapters/api/json"
	p "github.com/lejeunel/go-image-annotator/adapters/api/json/user"
	"github.com/lejeunel/go-image-annotator/adapters/api/models"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	"github.com/lejeunel/go-image-annotator/use-cases/user/create"
	upd "github.com/lejeunel/go-image-annotator/use-cases/user/update-privileges"
)

func (s *Server) CreateUser(w http.ResponseWriter, r *http.Request) {
	body, ok := json.MustDecodeJSON[models.NewUser](w, r)
	if !ok {
		return
	}

	req := create.Request{Id: body.Id}
	if body.Roles != nil {
		req.Roles = *body.Roles
	}
	if body.Groups != nil {
		req.Groups = *body.Groups
	}

	s.User.Create.Execute(
		r.Context(), req, p.NewPresenter(w, s.Logger))
}

func (s *Server) DeleteUserById(w http.ResponseWriter, r *http.Request, id string) {
	s.User.Delete.Execute(r.Context(), id, p.NewPresenter(w, s.Logger))
}

func (s *Server) WhoAmI(w http.ResponseWriter, r *http.Request) {
	user := u.IdentityFromContext(r.Context())

	if user != nil {
		json.WriteJSON(w, 200, User{
			Id:     user.Id,
			Groups: user.GroupNames(),
			Roles:  user.RoleNames(),
		})
		return
	}
	http.Error(w, "failed fetching user's identity", http.StatusBadRequest)
}

func (s *Server) UpdateUser(w http.ResponseWriter, r *http.Request, id string) {
	body, ok := json.MustDecodeJSON[models.UserPrivileges](w, r)
	if !ok {
		return
	}

	req := upd.Request{Id: id, Groups: body.Groups, Roles: body.Roles}
	s.User.UpdatePrivileges.Execute(r.Context(), req, p.NewPresenter(w, s.Logger))
}
