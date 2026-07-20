package user

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	cmp "github.com/lejeunel/go-image-annotator/adapters/web/components"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	rt "github.com/lejeunel/go-image-annotator/routes"
	rat "github.com/lejeunel/go-image-annotator/use-cases/user/renew-access-token"
)

type Server struct {
	b.UserDashboardBuilder
	b.AdminPageBuilder
	RenewAPITokenItr rat.Interactor
}

func (s *Server) Admin(w http.ResponseWriter, r *http.Request) {
	pb := s.AdminPageBuilder
	pb.SetUserIdentity(r.Context())
	pb.Build().Render(w)

}

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
	s.RenewAPITokenItr.Execute(r.Context(),
		user.Id, cmp.NewAPITokenPresenter(w))
}
func (s *Server) Route(r chi.Router, mws ...func(http.Handler) http.Handler) {

	r.Group(func(r chi.Router) {
		r.Use(mws...)
		r.Get(rt.UserDashboard, s.UserDashboard)
		r.Get(rt.Admin, s.Admin)
		r.Get(rt.NewAPIToken, s.NewAPIToken)
	})
}

func New(pb b.PageBuilder, i rat.Interactor) Server {
	return Server{b.NewUserDashboardBuilder(pb), b.NewAdminPageBuilder(pb), i}
}
