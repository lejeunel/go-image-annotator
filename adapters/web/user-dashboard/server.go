package user

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	cmp "github.com/lejeunel/go-image-annotator/adapters/web/components"
	"github.com/lejeunel/go-image-annotator/adapters/web/htmx"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	rt "github.com/lejeunel/go-image-annotator/routes"
	cpw "github.com/lejeunel/go-image-annotator/use-cases/user/change-password"
	rat "github.com/lejeunel/go-image-annotator/use-cases/user/renew-access-token"
)

type Server struct {
	b.UserDashboardBuilder
	RenewAPITokenItr  rat.Interactor
	ChangePasswordItr cpw.Interactor
}

func New(pb b.PageBuilder, i rat.Interactor, c cpw.Interactor) Server {
	return Server{b.NewUserDashboardBuilder(pb),
		i, c}
}

func (s *Server) UserDashboard(w http.ResponseWriter, r *http.Request) {
	udb := s.UserDashboardBuilder
	udb.SetUserIdentity(r.Context())
	udb.SetActiveSection(cmp.NoPageActive)
	udb.SetTitle("User Dashboard").SetHTMLTitle("Dashboard")
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
func (s *Server) ChangePassword(w http.ResponseWriter, r *http.Request) {
	user := u.IdentityFromContext(r.Context())
	if user == nil {
		http.Error(w, "failed getting user identity", http.StatusForbidden)
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form data", http.StatusBadRequest)
		return
	}
	s.ChangePasswordItr.Execute(r.Context(), cpw.Request{Id: user.Id, CurrentPassword: r.FormValue("password-current"),
		FirstPassword: r.FormValue("password"), SecondPassword: r.FormValue("password-repeat")},
		NewChangePasswordPresenter(w))

}
func (s *Server) Route(r chi.Router, mws ...func(http.Handler) http.Handler) {

	r.Group(func(r chi.Router) {
		r.Use(mws...)
		r.Get(rt.UserDashboard, s.UserDashboard)
		r.Get(rt.NewAPIToken, s.NewAPIToken)
		r.Post(rt.ChangePassword, s.ChangePassword)
	})
}

type ChangePasswordPresenter struct {
	writer http.ResponseWriter
	task   string
	htmx.ErrorPresenter
}

func NewChangePasswordPresenter(w http.ResponseWriter) ChangePasswordPresenter {
	task := "Change password"
	return ChangePasswordPresenter{w, task, htmx.NewErrorPresenter(task, w)}
}
func (p ChangePasswordPresenter) Success() {
	payload, _ := htmx.NotifySuccessPayloadAndReload(p.task, "Successfully changed password")
	p.writer.Header().Set("HX-Trigger", string(payload))
	p.writer.WriteHeader(http.StatusOK)
}
