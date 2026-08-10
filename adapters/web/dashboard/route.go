package dashboard

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	rt "github.com/lejeunel/go-image-annotator/routes"
)

func (s *Server) Route(r chi.Router, mws ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(mws...)
		r.Get(rt.DashboardUrl, s.Profile)
		r.Get(CredentialsUrl, s.Credentials)
		r.Get(rt.ListTasksUrl, s.ListTasks)
		r.Get(TaskRowUrl, s.TaskRow)
		r.Get(TaskDetailsUrl, s.TaskDetails)
		r.Get(NewAPITokenUrl, s.NewAPIToken)
		r.Post(ChangePasswordUrl, s.ChangePassword)
	})
}
