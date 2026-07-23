package user

import (
	"github.com/go-chi/chi/v5"
	rt "github.com/lejeunel/go-image-annotator/routes"
	"net/http"
)

var NewAPIToken = "/ui/new-api-token"
var ChangePassword = "/change-password"

func (s *Server) Route(r chi.Router, mws ...func(http.Handler) http.Handler) {

	r.Group(func(r chi.Router) {
		r.Use(mws...)
		r.Get(rt.UserDashboard, s.UserDashboard)
		r.Get(NewAPIToken, s.NewAPIToken)
		r.Post(ChangePassword, s.ChangePassword)
	})
}
