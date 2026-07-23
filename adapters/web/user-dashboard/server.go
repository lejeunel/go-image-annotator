package user

import (
	"net/http"

	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	cmp "github.com/lejeunel/go-image-annotator/adapters/web/components"
	"github.com/lejeunel/go-image-annotator/adapters/web/htmx"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	cpw "github.com/lejeunel/go-image-annotator/use-cases/user/change-password"
	rat "github.com/lejeunel/go-image-annotator/use-cases/user/renew-access-token"
)

type Server struct {
	b.PageBuilder
	RenewAPITokenItr  rat.Interactor
	ChangePasswordItr cpw.Interactor
}

func New(pb b.PageBuilder, i rat.Interactor, c cpw.Interactor) Server {
	return Server{pb, i, c}
}

func (s *Server) UserDashboard(w http.ResponseWriter, r *http.Request) {
	s.PageBuilder.SetUserIdentity(r.Context())
	RenderDashboard(r.Context(), s.PageBuilder, w)
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
	htmx.NotifySuccessPayloadAndReload(p.writer, p.task, "Successfully changed password")
}
