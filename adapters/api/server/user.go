package server

import (
	"net/http"

	"github.com/lejeunel/go-image-annotator/adapters/api/json"
	p "github.com/lejeunel/go-image-annotator/adapters/api/json/user"
	"github.com/lejeunel/go-image-annotator/adapters/api/models"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	"github.com/lejeunel/go-image-annotator/shared"
	"github.com/lejeunel/go-image-annotator/use-cases/user/create"
	upd "github.com/lejeunel/go-image-annotator/use-cases/user/update-privileges"
)

func (s *Server) CreateUser(w http.ResponseWriter, r *http.Request) {
	body, ok := json.MustDecodeJSON[models.NewUser](w, r)
	if !ok {
		return
	}
	s.User.Create.Execute(
		r.Context(),
		create.Request{
			Id: body.Id, Roles: body.Roles,
			Groups: body.Groups,
		}, p.NewPresenter(w, s.Logger))
}

func (s *Server) DeleteUser(w http.ResponseWriter, r *http.Request, id string) {
	s.User.Delete.Execute(r.Context(), id, p.NewPresenter(w, s.Logger))
}

func (s *Server) WhoAmI(w http.ResponseWriter, r *http.Request) {
	user := u.IdentityFromContext(r.Context())

	if user != nil {
		shared.WriteJSON(w, 200, User{
			Id:     user.Id,
			Groups: user.GroupNames(),
			Roles:  user.RoleNames(),
		})
		return
	}

	json.WriteError(w, http.StatusBadRequest, "failed fetching user's identity")
}

func (s *Server) UpdateUser(w http.ResponseWriter, r *http.Request, id string) {
	body, ok := json.MustDecodeJSON[models.UserPrivileges](w, r)
	if !ok {
		return
	}

	req := upd.Request{Id: id, Groups: body.Groups, Roles: body.Roles}
	s.User.UpdatePrivileges.Execute(r.Context(), req, p.NewPresenter(w, s.Logger))
}
