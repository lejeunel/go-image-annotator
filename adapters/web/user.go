package web

import (
	"net/http"

	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	p "github.com/lejeunel/go-image-annotator/shared/identity_provider"
	rt "github.com/lejeunel/go-image-annotator/use-cases/user/renew-access-token"
)

func (s *Server) UserDashboard(w http.ResponseWriter, r *http.Request) {
	p := s.PageBuilder
	p.SetUserIdentityFromContext(r.Context())
	p.SetActive(b.NoPageActive)
	udb := s.UserDashboardBuilder.SetUserIdentityFromContext(r.Context())
	p.SetTitle("User Dashboard")
	p.SetContent(udb.Build())
	p.Render(w)

}
func (s *Server) NewAPIToken(w http.ResponseWriter, r *http.Request) {
	user := p.IdentityFromContext(r.Context())
	if user == nil {
		http.Error(w, "failed getting user identity", http.StatusForbidden)
	}
	s.Interactors.User.RenewToken.Execute(r.Context(),
		rt.Request{Id: user.Id}, b.NewAPITokenPresenter(w))
}
