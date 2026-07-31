package dashboard

import (
	"github.com/go-chi/chi/v5"
	rt "github.com/lejeunel/go-image-annotator/routes"
	"net/http"
)

func (s *Server) Route(r chi.Router, mws ...func(http.Handler) http.Handler) {

	r.Group(func(r chi.Router) {
		r.Use(mws...)
		r.Get(rt.DashboardUrl, s.Credentials)
		r.Get(CredentialsUrl, s.Credentials)
		r.Get(LogsUrl, s.ListTasks)
		r.Get(NewAPITokenUrl, s.NewAPIToken)
		r.Post(ChangePasswordUrl, s.ChangePassword)
	})
}
