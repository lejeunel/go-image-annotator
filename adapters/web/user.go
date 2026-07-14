package web

import (
	"net/http"

	cmp "github.com/lejeunel/go-image-annotator/adapters/web/components"
	u "github.com/lejeunel/go-image-annotator/entities/user"
)

func (s *Server) UserDashboard(w http.ResponseWriter, r *http.Request) {
	udb := s.UserDashboardBuilder
	udb.SetUserIdentity(r.Context())
	udb.SetActive(cmp.NoPageActive)
	udb.SetTitle("User Dashboard")
	udb.Build().Render(w)

}
func (s *Server) NewAPIToken(w http.ResponseWriter, r *http.Request) {
	user := u.IdentityFromContext(r.Context())
	if user == nil {
		http.Error(w, "failed getting user identity", http.StatusForbidden)
	}
	s.Interactors.User.RenewToken.Execute(r.Context(),
		user.Id, cmp.NewAPITokenPresenter(w))
}
