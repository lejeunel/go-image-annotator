package web

import (
	"net/http"

	cmp "github.com/lejeunel/go-image-annotator/adapters/web/components"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	rt "github.com/lejeunel/go-image-annotator/use-cases/user/renew-access-token"
)

func (s *Server) UserDashboard(w http.ResponseWriter, r *http.Request) {
	p := s.PageBuilder
	p.SetUserIdentity(r.Context())
	p.SetActive(cmp.NoPageActive)
	udb := s.UserDashboardBuilder.SetUserIdentityFromContext(r.Context())
	p.SetTitle("User Dashboard")
	p.SetContent(udb.Build(), nil)
	p.Render(w)

}
func (s *Server) NewAPIToken(w http.ResponseWriter, r *http.Request) {
	user := u.IdentityFromContext(r.Context())
	if user == nil {
		http.Error(w, "failed getting user identity", http.StatusForbidden)
	}
	s.Interactors.User.RenewToken.Execute(r.Context(),
		rt.Request{Id: user.Id}, cmp.NewAPITokenPresenter(w))
}
